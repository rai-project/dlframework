package client

import "github.com/spf13/cobra"

var datasetCmd = &cobra.Command{
	Use:   "dataset",
	Short: "Request MLModelScope agents to predict a dataset",
	RunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
