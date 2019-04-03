package server

import (
	dllayer "github.com/rai-project/dllayer/cmd"
	evalcmd "github.com/rai-project/evaluation/cmd"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information from MLModelScope",
}

func init() {
	infoCmd.AddCommand(dllayer.FlopsInfoCmd)
	infoCmd.AddCommand(evalcmd.EvaluationCmd)
	infoCmd.AddCommand(infoModelsCmd)
}
