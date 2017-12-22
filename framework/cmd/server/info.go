package server

import (
	"fmt"

	dllayer "github.com/rai-project/dllayer/cmd"
	evaluations "github.com/rai-project/evaluation/cmd"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get model and other information from CarML",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside infoCmd Run with framework %v\n", framework)
	},
}

func init() {
	infoCmd.AddCommand(dllayer.FlopsInfoCmd)
	infoCmd.AddCommand(evaluations.EvaluationCmd)
}
