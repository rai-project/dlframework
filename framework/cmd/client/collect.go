package client

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:     "collectCmd",
	Short:   "collectCmd",
	Aliases: []string{"collect"},
	Long:    `collectCmd`,
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("urlsfile, batch size upperbound and output dir need to be provided")
		}
		batchUpper, _ := strconv.Atoi(args[1])
		for ii := 1; ii < batchUpper; ii = ii * 2 {
			err := urlsCmd.RunE(c, []string{args[0], strconv.Itoa(ii), args[2]})
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(collectCmd)
}
