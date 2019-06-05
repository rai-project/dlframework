package server

import (
	dllayer "github.com/rai-project/dllayer/cmd"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about the built-in models",
}

func init() {
	infoCmd.AddCommand(dllayer.FlopsInfoCmd)
	infoCmd.AddCommand(infoModelsCmd)
}
