package server

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/evaluation"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/tracer"
	"github.com/rai-project/tracer/jaeger"
	tracerutils "github.com/rai-project/tracer/utils"
	"github.com/rai-project/tracer/zipkin"
	"github.com/spf13/cobra"
)

var (
	modelName             string
	modelVersion          string
	useGPU                bool
	tracePreprocess       bool
	batchSize             int
	partitionListSize     int
	publishToDatabase     bool
	publishPredictions    bool
	failOnFirstError      bool
	traceLevelName        string
	traceLevel            tracer.Level = tracer.MODEL_TRACE
	tracerAddress         string
	databaseAddress       string
	databaseName          string
	databaseEndpoints     []string
	db                    database.Database
	evaluationTable       *evaluation.EvaluationCollection
	modelAccuracyTable    *evaluation.ModelAccuracyCollection
	performanceTable      *evaluation.PerformanceCollection
	inputPredictionsTable *evaluation.InputPredictionCollection
	DefaultChannelBuffer  = 100000
	fixTracerEndpoints    = tracerutils.FixEndpoints("http://", "9411", "/api/v1/spans")
	baseDir               string
)

var predictCmd = &cobra.Command{
	Use:   "predict",
	Short: "Predict using the agent",
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		rootCmd := c.Parent()
		for rootCmd.HasParent() {
			rootCmd = rootCmd.Parent()
		}
		rootCmd.PersistentPreRunE(c, args)

		if partitionListSize == 0 {
			partitionListSize = batchSize
		}
		traceLevel = tracer.LevelFromName(traceLevelName)
		if databaseName == "" {
			databaseName = config.App.Name
		}
		if databaseName != "" {
			databaseName = strings.Replace(databaseName, ".", "_", -1)
		}
		if databaseAddress == "" {
			if len(mongodb.Config.Endpoints) == 0 {
				panic("no database enpoint found")
			}
			databaseAddress = mongodb.Config.Endpoints[0]
		}
		if databaseAddress != "" {
			databaseEndpoints = []string{databaseAddress}
		}
		if tracerAddress != "" {
			tracerHost := getTracerHostAddress(tracerAddress)
			jaeger.Config.Endpoints = fixTracerEndpoints([]string{tracerHost})
			zipkin.Config.Endpoints = fixTracerEndpoints([]string{tracerHost})
			tracer.ResetStd()
		} else {
			tracerAddress = getTracerHostAddress(jaeger.Config.Endpoints[0])
		}
		if useGPU && !nvidiasmi.HasGPU {
			return errors.New("unable to find gpu on the system")
		}
		return nil
	},
}

func init() {
	predictCmd.PersistentFlags().StringVar(&modelName, "model_name", "MobileNet_v1_1.0_224", "the name of the model to use for prediction")
	predictCmd.PersistentFlags().StringVar(&modelVersion, "model_version", "1.0", "the version of the model to use for prediction")
	predictCmd.PersistentFlags().IntVarP(&batchSize, "batch_size", "b", 1, "the batch size to use while performing inference")
	predictCmd.PersistentFlags().IntVarP(&partitionListSize, "partition_list_size", "p", 0, "the chunk size to partition the input list. By default this is the same as the batch size")
	predictCmd.PersistentFlags().BoolVar(&useGPU, "use_gpu", false, "whether to enable the gpu. An error is returned if the gpu is not available")
	predictCmd.PersistentFlags().BoolVar(&tracePreprocess, "trace_preprocess", true, "whether to trace the preproessing steps. By defatult it is set to true")
	predictCmd.PersistentFlags().BoolVar(&failOnFirstError, "fail_on_error", false, "turning on causes the process to terminate/exit upon first inference error. This is useful since some inferences will result in an error because they run out of memory")
	predictCmd.PersistentFlags().BoolVar(&publishToDatabase, "publish", false, "whether to publish the evaluation to database. Turning this off will not publish anything to the database. This is ideal for using carml within profiling tools or performing experiments where the terminal output is sufficient.")
	predictCmd.PersistentFlags().BoolVar(&publishPredictions, "publish_predictions", false, "whether to publish prediction results to database. This will store all the probability outputs for the evaluation in the database which could be a few gigabytes of data for one dataset")
	predictCmd.PersistentFlags().StringVar(&traceLevelName, "trace_level", traceLevel.String(), "the trace level to use while performing evaluations")
	predictCmd.PersistentFlags().StringVar(&tracerAddress, "tracer_address", "", "the address of the jaeger or the zipking trace server")
	predictCmd.PersistentFlags().StringVar(&databaseName, "database_name", "", "the name of the database to publish the evaluation results to. By default the app name in the config `app.name` is used")
	predictCmd.PersistentFlags().StringVar(&databaseAddress, "database_address", "", "the address of the mongo database to store the results. By default the address in the config `database.endpoints` is used")
	predictCmd.PersistentFlags().StringVar(&baseDir, "base_dir", "results", "the folder path to store the results. By default 'results' is used")

	predictCmd.AddCommand(predictDatasetCmd)
	predictCmd.AddCommand(predictUrlsCmd)
	// predictCmd.AddCommand(predictWorkloadCmd)
	// predictCmd.AddCommand(predictQPSCmd)
}
