package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rai-project/batching"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/options"
	common "github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/dlframework/steps"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/synthetic_load"
	"github.com/rai-project/tracer"
	"github.com/rai-project/uuid"
	"github.com/spf13/cobra"
	"github.com/ulule/deepcopier"
	pb "gopkg.in/cheggaaa/pb.v2"
)

var predictWorkloadCmd = &cobra.Command{
	Use:     "workload",
	Short:   "Evaluates the workload using the specified model and framework",
	Aliases: []string{"work-load"},
	RunE: func(c *cobra.Command, args []string) error {
		span, ctx := tracer.StartSpanFromContext(context.Background(), traceLevel, "workload")
		defer func() {
			if span != nil {
				span.Finish()
			}
		}()

		predictorFramework, err := agent.GetPredictor(framework)
		if err != nil {
			return errors.Wrapf(err,
				"⚠️ failed to get predictor for %s. make sure you have "+
					"imported the framework's predictor package",
				framework.MustCanonicalName(),
			)
		}

		model, err := framework.FindModel(modelName + ":" + modelVersion)
		if err != nil {
			return err
		}

		var dc map[string]int32
		if useGPU {
			if !nvidiasmi.HasGPU {
				return errors.New("not gpu found")
			}
			dc = map[string]int32{"GPU": 0}
		} else {
			dc = map[string]int32{"CPU": 0}
		}

		execOpts := &dl.ExecutionOptions{
			TraceLevel: dl.ExecutionOptions_TraceLevel(
				dl.ExecutionOptions_TraceLevel_value[traceLevel.String()],
			),
			DeviceCount: dc,
		}
		predOpts := &dl.PredictionOptions{
			FeatureLimit:     10,
			BatchSize:        int32(batchSize),
			ExecutionOptions: execOpts,
		}

		predictor, err := predictorFramework.Load(ctx, *model, options.PredictorOptions(predOpts))
		if err != nil {
			return err
		}

		preprocessOptions, err := predictor.GetPreprocessOptions(nil) // disable tracing
		if err != nil {
			return err
		}

		var imagePredictor common.ImagePredictor

		err = deepcopier.Copy(predictor).To(&imagePredictor)
		if err != nil {
			return errors.Errorf("failed to copy to an image predictor for %v", model.MustCanonicalName())
		}

		var bar *pb.ProgressBar

		println("Starting workload generation process")

		batchQueue := make(chan []byte)
		go func() {
			defer close(batchQueue)
			opts := []synthetic_load.Option{synthetic_load.Context(ctx),
				synthetic_load.QPS(512),
				synthetic_load.InputGenerator(func(idx int) ([]byte, error) {
					return []byte("http://ww4.hdnux.com/photos/41/15/35/8705883/4/920x920.jpg"), nil
				}),
				synthetic_load.InputRunner(batchingRunner{
					queue: batchQueue,
				}),
			}
			tr := synthetic_load.NewTrace(opts...)
			bar = newProgress("workload", len(tr))
			latency := tr.Replay(opts...)
			fmt.Printf("qps = %v latency = %v", tr.QPS(), latency)
		}()

		partlabels := map[string]string{}

		batching.NewNaive(
			func(data [][]byte) {
				defer bar.Add(len(data))
				input := make(chan interface{}, DefaultChannelBuffer)
				go func() {
					defer close(input)
					for _, url := range data {
						id := uuid.NewV4()
						lbl := steps.NewIDWrapper(id, string(url))
						partlabels[lbl.GetID()] = "" // no label for the input url
						input <- lbl
					}
				}()

				output := pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
					Then(steps.NewReadURL()).
					Then(steps.NewReadImage(preprocessOptions)).
					Then(steps.NewPreprocessImage(preprocessOptions)).
					Run(input)
				var images []interface{}
				for out := range output {
					images = append(images, out)
				}

				input = make(chan interface{}, DefaultChannelBuffer)
				go func() {
					input <- images
				}()
				output = pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
					Then(steps.NewPredictImage(predictor)).
					Run(input)

			},
			batchQueue,
			batching.BatchSize(batchSize),
		)
		bar.Finish()
		return nil
	},
}

type batchingRunner struct {
	queue chan []byte
}

func (s batchingRunner) Run(input []byte, onFinish func()) error {
	s.queue <- input
	go func() { // HACK
		for {
			if len(s.queue) == cap(s.queue) {
				onFinish()
			}
		}
	}()
	return nil
}
