package server

import (
	"os"
	"path/filepath"
	"strings"

	"encoding/csv"
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
	fullFlops         bool
	humanFlops        bool
	flopsOutputFormat string
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

		tableWriter := tablewriter.NewWriter(os.Stdout)
		csvWriter := csv.NewWriter(os.Stdout)

		flopsOutputFormat := strings.ToLower(flopsOutputFormat)

		writeHeader := func(header []string) {
			switch flopsOutputFormat {
			case "table":
				tableWriter.SetHeader(header)
			case "csv":
				csvWriter.Write(header)
			}
		}

		writeRecord := func(row []string) {
			switch flopsOutputFormat {
			case "table":
				tableWriter.Append(row)
			case "csv":
				csvWriter.Write(row)
			}
		}

		flush := func() {
			switch flopsOutputFormat {
			case "table":
				tableWriter.Render()
			case "csv":
				csvWriter.Flush()
			}
		}

		defer flush()

		if fullFlops {

			infos := net.LayerInformations()

			writeHeader([]string{"LayerName", "LayerType", "#MultiplyAdds", "#Additions", "#Divisions", "#Exponentiations", "#Comparisons", "#General"})

			for _, info := range infos {
				flops := info.Flops()
				writeRecord([]string{
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
			return nil
		}

		info := net.FlopsInformation()

		writeHeader([]string{"Flop Type", "#"})
		writeRecord([]string{"MultipleAdds", flopsToString(info.MultiplyAdds)})
		writeRecord([]string{"Additions", flopsToString(info.Additions)})
		writeRecord([]string{"Divisions", flopsToString(info.Divisions)})
		writeRecord([]string{"Exponentiations", flopsToString(info.Exponentiations)})
		writeRecord([]string{"Comparisons", flopsToString(info.Comparisons)})
		writeRecord([]string{"General", flopsToString(info.General)})

		return nil
	},
}

func init() {
	flopsInfoCmd.PersistentFlags().StringVar(&modelName, "model_name", "BVLC-AlexNet", "modelName")
	flopsInfoCmd.PersistentFlags().StringVar(&modelVersion, "model_version", "1.0", "modelVersion")
	flopsInfoCmd.PersistentFlags().BoolVar(&humanFlops, "human", false, "print flops in human form")
	flopsInfoCmd.PersistentFlags().BoolVar(&fullFlops, "full", false, "print all information about flops")
	flopsInfoCmd.PersistentFlags().StringVarP(&flopsOutputFormat, "format", "f", "table", "print format to use")
}
