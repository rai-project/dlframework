package client

import "github.com/spf13/cobra"

var imagesCmd = &cobra.Command{
	Use:   "imagesCmd",
	Short: "imagesCmd",
	Long:  `imagesCmd`,
	RunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
