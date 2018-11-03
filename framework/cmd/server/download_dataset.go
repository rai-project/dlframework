package server

import (
	"context"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/rai-project/dldataset"
	"github.com/spf13/cobra"
)

var (
	downloadDatasetCategory string
	downloadDatasetName     string
	downloadAll             bool
)

func datasetDownload(ctx context.Context, datasetCategory, datasetName string) error {
	dataset, err := dldataset.Get(datasetCategory, datasetName)
	if err != nil {
		return err
	}
	defer dataset.Close()

	err = dataset.Download(ctx)
	if err != nil {
		return err
	}

	return nil
}

func datasetDownloadAll(ctx context.Context) error {
	datasets := dldataset.Datasets()
	for _, dataset := range datasets {
		splt := strings.Split(dataset, "/")
		category, name := splt[0], splt[1]
		pp.Println(splt)
		err := datasetDownload(ctx, category, name)
		if err != nil {
			log.WithError(err).Errorf("failed to download %v", dataset)
			return err
		}
	}

	return nil
}

var downloadDatasetCmd = &cobra.Command{
	Use:   "dataset",
	Short: "Download MLModelScope datasets",
	RunE: func(c *cobra.Command, args []string) error {
		ctx := context.Background()

		if downloadAll == true {
			err := datasetDownloadAll(ctx)
			if err != nil {
				return err
			}
		}

		err := datasetDownload(ctx, downloadDatasetCategory, downloadDatasetName)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	downloadDatasetCmd.PersistentFlags().StringVar(&downloadDatasetCategory, "category", "vision", "dataset category (e.g. \"vision\")")
	downloadDatasetCmd.PersistentFlags().StringVar(&downloadDatasetName, "name", "ilsvrc2012_validation_224", "dataset name (e.g. \"ilsvrc2012_validation_folder\")")
	downloadDatasetCmd.PersistentFlags().BoolVar(&downloadAll, "all", false, "download all the available datasets")
}
