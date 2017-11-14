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
	"github.com/rai-project/utils"
	"github.com/spf13/cobra"
)

var (
	fullFlops  bool
	humanFlops bool
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

		flopsToString := func(e int64) string {
			return fmt.Sprintf("%v", e)
		}
		if humanFlops {
			flopsToString = func(e int64) string {
				return utils.Flops(uint64(e))
			}
		}

		if fullFlops {
			infos := net.LayerInformations()

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"LayerName", "LayerType", "#MultiplyAdds", "#Additions", "#Divisions", "#Exponentiations", "#Comparisons", "#General"})
			for _, info := range infos {
				flops := info.Flops()
				table.Append([]string{
					info.Name(),
					info.Type(),
					flopsToString(flops.MultiplyAdds),
					flopsToString(flops.Additions),
					flopsToString(flops.Divisions),
					flopsToString(flops.Exponentiations),
					flopsToString(flops.Comparisons),
					flopsToString(flops.General),
				})
			}
			table.Render()
			return nil
		}

		info := net.FlopsInformation()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Flop Type", "#"})
		table.Append([]string{"MultipleAdds", flopsToString(info.MultiplyAdds)})
		table.Append([]string{"Additions", flopsToString(info.Additions)})
		table.Append([]string{"Divisions", flopsToString(info.Divisions)})
		table.Append([]string{"Exponentiations", flopsToString(info.Exponentiations)})
		table.Append([]string{"Comparisons", flopsToString(info.Comparisons)})
		table.Append([]string{"General", flopsToString(info.General)})
		table.Render()

		return nil
	},
}

func init() {
	flopsInfoCmd.PersistentFlags().StringVar(&modelName, "model_name", "BVLC-AlexNet", "modelName")
	flopsInfoCmd.PersistentFlags().StringVar(&modelVersion, "model_version", "1.0", "modelVersion")
	flopsInfoCmd.PersistentFlags().BoolVar(&humanFlops, "human", false, "print flops in human form")
	flopsInfoCmd.PersistentFlags().BoolVar(&fullFlops, "full", false, "print all information about flops")
}
