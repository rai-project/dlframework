package server

import "github.com/spf13/cobra"

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads CarML resource files",
}

func init() {
	downloadCmd.AddCommand(downloadDatasetCmd)
	downloadCmd.AddCommand(downloadModelsCmd)
}
