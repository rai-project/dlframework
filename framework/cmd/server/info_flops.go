// +build ignore

package server

import (
	"os"

	dllayer "github.com/rai-project/dllayer/cmd"
	evalcmd "github.com/rai-project/evaluation/cmd"
	"github.com/spf13/cobra"
)

var infoMLArcCmd = &cobra.Command{
	Use: "mlarc",
	Aliases: []string{
		"ml-arc",
		"all",
		"mlarc-web",
	},
	Short: "Get mlarc information from MLModelScope",
	RunE: func(cmd *cobra.Command, args []string) error {
		dllayer.FlopsInfoCmd.SetArgs(os.Args[2:])
		evalcmd.EvaluationCmd.SetArgs(append([]string{"all"}, os.Args[2:]...))

		dllayer.FlopsInfoCmd.Execute()
		evalcmd.EvaluationCmd.Execute()

		return nil
	},
}
