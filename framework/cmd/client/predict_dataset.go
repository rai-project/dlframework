package client

import "github.com/spf13/cobra"

var datasetCmd = &cobra.Command{
	Use:   "datasetCmd",
	Short: "datasetCmd",
	Long:  `datasetCmd`,
	RunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
