package server

import "github.com/spf13/cobra"

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads carml resource files",
}

func init() {
	downloadCmd.AddCommand(downloadDatasetCmd)
}
