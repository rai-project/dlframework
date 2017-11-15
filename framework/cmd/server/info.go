package server

import "github.com/spf13/cobra"

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get model and other information from CarML",
}

func init() {
	infoCmd.AddCommand(flopsInfoCmd)
}
