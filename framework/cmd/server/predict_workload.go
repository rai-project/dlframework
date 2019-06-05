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
	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/batching"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/options"
	common "github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/dlframework/steps"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/synthetic_load"
	"github.com/rai-project/tracer"
	"github.com/schollz/progressbar"
	"github.com/spf13/cobra"
	"github.com/ulule/deepcopier"
)

var (
	qps                    float64
	latencyBound           int64
	latencyBoundPercentile float64
	minDuration            int64
	minQueries             int
	maxQpsSearchIterations int
	imagePath              string
)

func computeLatency(qps float64) (trace synthetic_load.Trace, latency time.Duration, err error) {
	span, ctx := tracer.StartSpanFromContext(context.Background(), traceLevel, "workload")
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()

	predictorFramework, err := agent.GetPredictor(framework)
	if err != nil {
		err = errors.Wrapf(err,
			"⚠️ failed to get predictor for %s. make sure you have "+
				"imported the framework's predictor package",
			framework.MustCanonicalName(),
		)
		return
	}

	model, err := framework.FindModel(modelName + ":" + modelVersion)
	if err != nil {
		return
	}

	var dc map[string]int32
	if useGPU {
		if !nvidiasmi.HasGPU {
			err = errors.New("not gpu found")
			return
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

	predictor, err := predictorFramework.Load(
		ctx,
		*model,
		// options.Context(ctx),
		options.PredictorOptions(predOpts),
		// options.DisableFrameworkAutoTuning(true),
	)
	if err != nil {
		return
	}

	preprocessOptions, err := predictor.GetPreprocessOptions()
	if err != nil {
		return
	}

	var imagePredictor common.ImagePredictor

	err = deepcopier.Copy(predictor).To(&imagePredictor)
	if err != nil {
		err = errors.Errorf("failed to copy to an image predictor for %v", model.MustCanonicalName())
		return
	}

	var bar *progressbar.ProgressBar
	useBar := false

	println("Starting inference workload generation process")

	batchQueue := make(chan steps.IDer)
	outputQueue := new(sync.Map)
	go func() {
		defer close(batchQueue)

		input, err := ioutil.ReadFile(imagePath)
		if err != nil {
			panic(err)
		}

		opts := []synthetic_load.Option{
			synthetic_load.Context(ctx),
			synthetic_load.QPS(qps),
			synthetic_load.LatencyBoundPercentile(latencyBoundPercentile),
			synthetic_load.MinQueries(minQueries),
			synthetic_load.MinDuration(time.Duration(minDuration * int64(time.Millisecond))),
			synthetic_load.InputGenerator(func(idx int) ([]byte, error) {
				return input, nil
			}),
			synthetic_load.InputRunner(batchingRunner{
				inputQueue:  batchQueue,
				outputQueue: outputQueue,
				batchSize:   batchSize,
			}),
		}

		trace = synthetic_load.NewTrace(opts...)

		if useBar {
			bar = progressbar.NewOptions(len(trace), progressbar.OptionSetRenderBlankState(true))
		}

		latency, err = trace.Replay(opts...)
		if err != nil {
			return
		}
		// qps := trace.QPS()
		// fmt.Printf("qps = %v latency = %v \n", qps, latency)
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
				defer close(input)
				input <- images
			}()
			output = pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
				Then(steps.NewPredict(predictor)).
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

	return
}

var predictWorkloadCmd = &cobra.Command{
	Use:     "workload",
	Short:   "Evaluate the workload using the specified model and framework",
	Aliases: []string{"work-load"},
	RunE: func(c *cobra.Command, args []string) error {
		tr, latency, err := computeLatency(qps)
		if err != nil {
			return err
		}

		fmt.Printf("qps = %v, latency = %v\n",
			tr.QPS(),
			latency,
		)
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

func init() {
	sourcePath := sourcepath.MustAbsoluteDir()
	defaultImagePath := filepath.Join(sourcePath, "..", "_fixtures", "chicken.jpg")
	if !com.IsFile(defaultImagePath) {
		defaultImagePath = ""
	}

	predictWorkloadCmd.PersistentFlags().Float64Var(&qps, "initial_qps", 16, "the initial QPS")
	predictWorkloadCmd.PersistentFlags().Float64Var(&latencyBoundPercentile, "percentile", 95, "the minimum percent of queries meeting the latency bound")
	predictWorkloadCmd.PersistentFlags().Int64Var(&minDuration, "min_duration", 100, "the minimum duration of the trace in ms")
	predictWorkloadCmd.PersistentFlags().IntVar(&minQueries, "min_queries", 512, "the minimum number of queries")
	predictUrlsCmd.PersistentFlags().StringVar(&imagePath, "image_path", defaultImagePath, "the path to the image to perform the evaluations on.")
}
