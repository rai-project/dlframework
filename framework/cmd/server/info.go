package server

import (
	dllayer "github.com/rai-project/dllayer/cmd"
	evaluations "github.com/rai-project/evaluation/cmd"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get model and other information from CarML",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		dllayer.Framework = framework
	},
}

func init() {
	infoCmd.AddCommand(dllayer.FlopsInfoCmd)
	infoCmd.AddCommand(evaluations.EvaluationCmd)
}
