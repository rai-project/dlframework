package server

import (
	"context"
	"strings"

	"github.com/rai-project/dldataset"
)

func datasetDownload(ctx context.Context, datasetCategory, datasetName string) error {
	dataset, err := dldataset.Get(datasetCategory, datasetName)
	if err != nil {
		return err
	}
	return dataset.Download(ctx)
}

func datasetDownloadAll(ctx context.Context) error {
	datasets := dldataset.Datasets()
	for _, dataset := range datasets {
		splt := strings.Split(dataset, "/")
		category, name := splt[0], splt[1]
		err := datasetDownload(ctx, category, name)
		if err != nil {
			log.WithError(err).Errorf("failed to download %v", dataset)
		}
	}
	return nil
}
