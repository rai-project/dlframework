package server

import (
	"context"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/rai-project/dldataset"
	"github.com/spf13/cobra"
)

var (
	modelName    string
	modelVersion string
	all          bool
)

func modelDownload(ctx context.Context, datasetCategory, datasetName string) error {
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

func modelDownloadAll(ctx context.Context) error {
	datasets := dldataset.Datasets()
	for _, dataset := range datasets {
		splt := strings.Split(dataset, "/")
		category, name := splt[0], splt[1]
		pp.Println(splt)
		err := modelDownload(ctx, category, name)
		if err != nil {
			log.WithError(err).Errorf("failed to download %v", dataset)
			return err
		}
	}

	return nil
}

var downloadModelCmd = &cobra.Command{
	Use:   "downloadModel",
	Short: "download models",
	RunE: func(c *cobra.Command, args []string) error {
		ctx := context.Background()

		if downloadAll == true {
			err := modelDownloadAll(ctx)
			if err != nil {
				return err
			}
		}

		err := modelDownload(ctx, modelName, modelVersion)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	downloadModelCmd.PersistentFlags().StringVar(&modelName, "dataset_category", "vision", "dataset category (e.g. \"vision\")")
	downloadModelCmd.PersistentFlags().StringVar(&modelVersion, "dataset_name", "ilsvrc2012_validation_224", "dataset name (e.g. \"ilsvrc2012_validation_folder\")")
	downloadModelCmd.PersistentFlags().BoolVar(&all, "dowloadAll", false, "download all the available datasets")
}
