package server

import (
	"context"

	"github.com/k0kubun/pp"

	"github.com/pkg/errors"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/options"
	"github.com/rai-project/dlframework/steps"
	"github.com/rai-project/image/types"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/pipeline"
	"github.com/spf13/cobra"
)

var (
	datasetCategory      string
	datasetName          string
	modelName            string
	modelVersion         string
	batchSize            int
	partitionDatasetSize int
)

var (
	DefaultChannelBuffer = 1000
)

var datasetCmd = &cobra.Command{
	Use:   "dataset",
	Short: "dataset",
	RunE: func(c *cobra.Command, args []string) error {
		dataset, err := dldataset.Get(datasetCategory, datasetName)
		if err != nil {
			return err
		}
		defer dataset.Close()

		ctx := context.Background()

		err = dataset.Download(ctx)
		if err != nil {
			return err
		}

		fileList, err := dataset.List(ctx)
		if err != nil {
			return err
		}

		predictorFramework, err := agent.GetPredictor(framework)
		if err != nil {
			return errors.Wrapf(err,
				"⚠️ failed to get predictor for %s. make sure you have "+
					"imported the framework's predictor package",
				framework.MustCanonicalName(),
			)
		}

		model, err := framework.FindModel(modelName + "/" + modelVersion)
		if err != nil {
			return err
		}

		var dc map[string]int32
		if nvidiasmi.HasGPU {
			dc = map[string]int32{"GPU": 0}
		} else {
			dc = map[string]int32{"CPU": 0}
		}
		execOpts := &dlframework.ExecutionOptions{
			TraceLevel: dlframework.ExecutionOptions_TraceLevel(
				dlframework.ExecutionOptions_TraceLevel_value["NO_TRACE"]),
			DeviceCount: dc,
		}

		predOpts := &dlframework.PredictionOptions{
			FeatureLimit:     5,
			BatchSize:        uint32(batchSize),
			ExecutionOptions: execOpts,
		}

		predictor, err := predictorFramework.Load(ctx, *model, options.PredictorOptions(predOpts))
		if err != nil {
			return err
		}

		preprocessOptions, err := predictor.GetPreprocessOptions(ctx)
		if err != nil {
			return err
		}

		fileNameParts := partitionDataset(fileList, partitionDatasetSize)

		for _, part := range fileNameParts {
			partData := make([]*types.RGBImage, len(part))
			partlabels := make([]string, len(part))
			for ii, fileName := range part {
				lda, err := dataset.Get(ctx, fileName)
				if err != nil {
					return err
				}
				data, err := lda.Data()
				if err != nil {
					return err
				}
				rgbData := data.(*types.RGBImage)
				partData[ii] = rgbData
				partlabels[ii] = lda.Label()
			}

			input := make(chan interface{}, DefaultChannelBuffer)
			go func() {
				defer close(input)
				for _, img := range partData {
					input <- img
				}
			}()

			output := pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
				Then(steps.NewReadImage(preprocessOptions)).
				Then(steps.NewPreprocessImage(preprocessOptions)).
				Run(input)

			var outputs []interface{}
			for out := range output {
				outputs = append(outputs, out)
			}

			parts := agent.Partition(outputs, batchSize)

			input = make(chan interface{}, DefaultChannelBuffer)
			go func() {
				defer close(input)
				for _, part := range parts {
					input <- part
				}
			}()

			output = pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
				Then(steps.NewPredictImage(predictor)).
				Run(input)

			pp.Println(output)
		}

		return nil
	},
}

func partitionDataset(in []string, partitionSize int) (out [][]string) {
	cnt := (len(in)-1)/partitionSize + 1
	for i := 0; i < cnt; i++ {
		start := i * partitionSize
		end := (i + 1) * partitionSize
		if end > len(in) {
			end = len(in)
		}
		part := in[start:end]
		out = append(out, part)
	}

	return out
}

func init() {
	datasetCmd.PersistentFlags().StringVar(&datasetName, "dataset_category", "vision", "dataset category (e.g. \"vision\")")
	datasetCmd.PersistentFlags().StringVar(&datasetName, "dataset_name", "ilsvrc2012_validation_folder", "dataset name (e.g. \"ilsvrc2012_validation_folder\")")
	datasetCmd.PersistentFlags().StringVar(&modelName, "modelName", "BVLC-AlexNet", "modelName")
	datasetCmd.PersistentFlags().StringVar(&modelVersion, "modelVersion", "1.0", "modelVersion")
	datasetCmd.PersistentFlags().IntVarP(&batchSize, "batchSize", "b", 1, "batch size")
	datasetCmd.PersistentFlags().IntVarP(&partitionDatasetSize, "partitionDatasetSize", "p", 32, "partition dataset size")
}
