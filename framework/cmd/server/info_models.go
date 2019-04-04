package server

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var infoModelsCmd = &cobra.Command{
	Use:     "models",
	Aliases: []string{},
	Short:   "Get the model names and version registered by the agent",
	Run: func(cmd *cobra.Command, args []string) {
		models := framework.Models()
		if len(models) == 0 {
			fmt.Println("No Models")
			return
		}

		tbl := tablewriter.NewWriter(os.Stdout)
		tbl.SetHeader([]string{"Name", "Version", "Cannonical Name"})
		for _, model := range models {
			tbl.Append([]string{
				model.Name,
				model.Version,
				model.MustCanonicalName(),
			})
		}
		tbl.Render()
	},
}
