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
	predictCmd.AddCommand(predictDatasetCmd)
	predictCmd.AddCommand(predictUrlsCmd)
}
