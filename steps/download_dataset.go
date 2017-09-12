package steps

import (
	"sync"

	"golang.org/x/net/context"

	"github.com/rai-project/dldataset"
	"github.com/rai-project/pipeline"
)

type downloadDataset struct {
	base
	sync.Mutex
	isDownloaded bool
	dataset      dldataset.Dataset
}

func NewDownloadDataset(dataset dldataset.Dataset) pipeline.Step {
	res := downloadDataset{
		dataset:      dataset,
		isDownloaded: false,
		base: base{
			info: "DownloadDataset",
		},
	}
	res.doer = res.do
	return res
}

func (p downloadDataset) do(ctx context.Context, in0 interface{}) interface{} {
	defer func() { p.isDownloaded = true }()
	p.Lock()
	defer p.Unlock()
	if p.isDownloaded {
		return in0
	}
	err := p.dataset.Download(ctx)
	if err != nil {
		return err
	}
	return in0
}

func (p downloadDataset) Close() error {
	return nil
}
