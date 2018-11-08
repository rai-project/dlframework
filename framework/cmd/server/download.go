package server

import "github.com/spf13/cobra"

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads MLModelScope resource files",
}

func init() {
	downloadCmd.AddCommand(downloadDatasetCmd)
	downloadCmd.AddCommand(downloadModelsCmd)
}
