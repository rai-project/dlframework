package server

import (
	//"os"

	//dllayer "github.com/rai-project/dllayer/cmd"
	evaluations "github.com/rai-project/evaluation/cmd"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get model and other information from CarML",
}

/*
var infoMLArcCmd = &cobra.Command{
	Use: "mlarc",
	Aliases: []string{
		"ml-arc",
		"all",
		"mlarc-web",
	},
	Short: "Get mlarc information from CarML",
	RunE: func(cmd *cobra.Command, args []string) error {
		dllayer.FlopsInfoCmd.SetArgs(os.Args[2:])
		evaluations.EvaluationCmd.SetArgs(append([]string{"all"}, os.Args[2:]...))

		dllayer.FlopsInfoCmd.Execute()
		evaluations.EvaluationCmd.Execute()

		return nil
	},
}
*/

func init() {
	//infoCmd.AddCommand(dllayer.FlopsInfoCmd)
	infoCmd.AddCommand(evaluations.EvaluationCmd)
	infoCmd.AddCommand(infoModelsCmd)
	//infoCmd.AddCommand(infoMLArcCmd)
}
