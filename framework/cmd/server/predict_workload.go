package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	sourcepath "github.com/GeertJohan/go-sourcepath"
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
		useBar := false

		println("Starting inference workload generation process")

		batchQueue := make(chan steps.IDer, DefaultChannelBuffer)
		outputQueue := new(sync.Map)
		go func() {
			defer close(batchQueue)

			imagePath := filepath.Join(sourcepath.MustAbsoluteDir(), "_fixtures", "chicken.jpg")
			input, err := ioutil.ReadFile(imagePath)
			if err != nil {
				panic(err)
			}
			opts := []synthetic_load.Option{
				synthetic_load.Context(ctx),
				synthetic_load.QPS(512),
				synthetic_load.MinQueries(64),
				synthetic_load.LatencyBoundPercentile(0.99),
				synthetic_load.MinDuration(1 * time.Second),
				synthetic_load.InputGenerator(func(idx int) ([]byte, error) {
					return input, nil
				}),
				synthetic_load.InputRunner(batchingRunner{
					inputQueue:  batchQueue,
					outputQueue: outputQueue,
					batchSize:   batchSize,
				}),
			}
			tr := synthetic_load.NewTrace(opts...)
			if useBar {
				bar = newProgress("inference workload prediction", len(tr))
			}
			latency := tr.Replay(opts...)
			qps := tr.QPS()
			fmt.Printf("qps = %v latency = %v \n", qps, latency)
		}()

		btch, err := batching.NewNaive(
			func(data []steps.IDer) {
				if useBar {
					defer bar.Add(len(data))
				}

				input := make(chan interface{}, DefaultChannelBuffer)
				go func() {
					defer close(input)
					for _, elem := range data {
						input <- elem
					}
				}()
				output := pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
					Then(steps.NewReadImage(preprocessOptions)).
					Then(steps.NewPreprocessImage(preprocessOptions)).
					Run(input)

				var images []interface{}
				for out := range output {
					images = append(images, out)
				}

				input = make(chan interface{})
				go func() {
					input <- images
				}()
				output = pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
					Then(steps.NewPredictImage(predictor)).
					Run(input)

				for out0 := range output {
					if err, ok := out0.(error); ok {
						panic(err)
					}

					out := out0.(steps.IDer)
					qu0, ok := outputQueue.Load(out.GetID())
					if !ok {
						panic("cannot find " + out.GetID() + " input output queue")
					}
					qu := qu0.(chan struct{})
					qu <- struct{}{}
				}
			},
			batchQueue,
			batching.BatchSize(batchSize),
		)
		if err != nil {
			panic(err)
		}
		btch.Wait()

		if useBar {
			bar.Finish()
		}

		return nil
	},
}

type workloadInput struct {
	id   string
	data io.Reader
}

func (w workloadInput) GetID() string {
	return w.id
}

func (w workloadInput) GetData() interface{} {
	return w.data
}

type batchingRunner struct {
	inputQueue  chan steps.IDer
	outputQueue *sync.Map
	batchSize   int
}

func (s batchingRunner) Run(tr synthetic_load.TraceEntry, bts []byte, onFinish func()) error {
	id := strconv.Itoa(tr.Index)
	s.inputQueue <- workloadInput{
		id:   id,
		data: bytes.NewBuffer(bts),
	}
	ch := make(chan struct{})
	s.outputQueue.Store(id, ch)

	go func() {
		for {
			select {
			case <-ch:
				onFinish()
				return
			}
		}
	}()
	return nil
}
