package server

import (
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/Unknwon/com"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dllayer/network"
	"github.com/spf13/cobra"
)

var (
	fullFlops bool
)

func cleanPath(path string) string {
	path = strings.Replace(path, ":", "_", -1)
	path = strings.Replace(path, " ", "_", -1)
	path = strings.Replace(path, "-", "_", -1)
	return strings.ToLower(path)
}

func getGraphPath(model *dlframework.ModelManifest) string {
	graphPath := filepath.Base(model.GetModel().GetGraphPath())
	wd, _ := model.WorkDir()
	return cleanPath(filepath.Join(wd, graphPath))
}

var flopsInfoCmd = &cobra.Command{
	Use:   "model",
	Short: "Get flops information about the model",
	RunE: func(c *cobra.Command, args []string) error {

		model, err := framework.FindModel(modelName + ":" + modelVersion)
		if err != nil {
			return err
		}

		graphPath := getGraphPath(model)
		if !com.IsFile(graphPath) {
			return errors.Errorf("file %v does not exist", graphPath)
		}

		net, err := network.NewCaffe(graphPath)
		if err != nil {
			return err
		}

		if fullFlops {
			infos := net.LayerInformations()

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"LayerName", "#MultiplyAdds", "#Additions", "#Divisions", "#Exponentiations", "#Comparisons", "#General"})
			for _, info := range infos {
				flops := info.Flops()
				table.Append([]string{
					"",
					fmt.Sprintf("%v", flops.MultiplyAdds),
					fmt.Sprintf("%v", flops.Additions),
					fmt.Sprintf("%v", flops.Divisions),
					fmt.Sprintf("%v", flops.Exponentiations),
					fmt.Sprintf("%v", flops.Comparisons),
					fmt.Sprintf("%v", flops.General),
				})
			}
			table.Render()
			return nil
		}

		info := net.FlopsInformation()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Flop Type", "#"})
		table.Append([]string{"MultipleAdds", fmt.Sprintf("%v", info.MultiplyAdds)})
		table.Append([]string{"Additions", fmt.Sprintf("%v", info.Additions)})
		table.Append([]string{"Divisions", fmt.Sprintf("%v", info.Divisions)})
		table.Append([]string{"Exponentiations", fmt.Sprintf("%v", info.Exponentiations)})
		table.Append([]string{"Comparisons", fmt.Sprintf("%v", info.Comparisons)})
		table.Append([]string{"General", fmt.Sprintf("%v", info.General)})
		table.Render()

		return nil
	},
}

func init() {
	flopsInfoCmd.PersistentFlags().StringVar(&modelName, "model_name", "BVLC-AlexNet", "modelName")
	flopsInfoCmd.PersistentFlags().StringVar(&modelVersion, "model_version", "1.0", "modelVersion")
	flopsInfoCmd.PersistentFlags().BoolVar(&fullFlops, "full", false, "print all information about flops")
}
