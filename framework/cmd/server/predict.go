package server

import (
	"time"

	"github.com/rai-project/tracer"
	"github.com/spf13/cobra"
	pb "gopkg.in/cheggaaa/pb.v2"
)

var (
	modelName            string
	modelVersion         string
	useGPU               bool
	batchSize            int
	publishEvaluation    bool
	publishPredictions   bool
	failOnFirstError     bool
	traceLevelName       string
	traceLevel           tracer.Level = tracer.APPLICATION_TRACE
	traceServerAddress   string
	databaseAddress      string
	databaseName         string
	databaseEndpoints    []string
	DefaultChannelBuffer = 100000
)

func newProgress(prefix string, count int) *pb.ProgressBar {
	// get the new original progress bar.
	//bar := pb.New(count).Prefix(prefix)
	// TODO: set prefix of bar
	bar := pb.New(count)
	//bar.Set("prefix", prefix)

	// Refresh rate for progress bar is set to 100 milliseconds.
	bar.SetRefreshRate(time.Millisecond * 100)

	bar.SetTemplateString(string(pb.Full))
	bar.Start()
	return bar
}

var predictCmd = &cobra.Command{
	Use:   "predict",
	Short: "Predicts using the MLModelScope agent",
}

func init() {
	predictCmd.PersistentFlags().StringVar(&modelName, "model_name", "BVLC-AlexNet", "the name of the model to use for prediction")
	predictCmd.PersistentFlags().StringVar(&modelVersion, "model_version", "1.0", "the version of the model to use for prediction")
	predictCmd.PersistentFlags().BoolVar(&useGPU, "gpu", false, "whether to enable the gpu. An error is returned if the gpu is not available")
	predictCmd.PersistentFlags().BoolVar(&failOnFirstError, "fail_on_error", false, "turning on causes the process to terminate/exit upon first inference error. This is useful since some inferences will result in an error because they run out of memory")
	predictCmd.PersistentFlags().BoolVar(&publishEvaluation, "publish", false, "whether to publish the evaluation to database. Turning this off will not publish anything to the database. This is ideal for using carml within profiling tools or performing experiments where the terminal output is sufficient.")
	predictCmd.PersistentFlags().BoolVar(&publishPredictions, "publish_predictions", false, "whether to publish prediction results to database. This will store all the probability outputs for the evaluation in the database which could be a few gigabytes of data for one dataset")
	predictCmd.PersistentFlags().StringVar(&traceLevelName, "trace_level", traceLevel.String(), "the trace level to use while performing evaluations")
	predictCmd.PersistentFlags().StringVar(&traceServerAddress, "tracer_address", "localhost:16686", "the address of the jaeger or the zipking trace server")
	predictCmd.PersistentFlags().StringVar(&databaseName, "database_name", "", "the name of the database to publish the evaluation results to. By default the app name in the config `app.name` is used")
	predictCmd.PersistentFlags().StringVar(&databaseAddress, "database_address", "", "the address of the mongo database to store the results. By default the address in the config `database.endpoints` is used")

	predictCmd.AddCommand(predictDatasetCmd)
	predictCmd.AddCommand(predictCmd)
}
