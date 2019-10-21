package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	tracePath  string
	outputPath string
)

var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Get various information from a trace",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running " + cmd.Name())
		return nil
	},
}

func init() {
	traceCmd.PersistentFlags().StringVar(&tracePath, "trace_path", "trace.json", "the path to the trace.")
	traceCmd.PersistentFlags().StringVarP(&outputPath, "output_path", "o", "", "the output path.")

	traceCmd.AddCommand(chromeTraceCmd)
	traceCmd.AddCommand(flameGraphCmd)
	// traceCmd.AddCommand(layersCmd)
}
