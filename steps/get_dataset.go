package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/dldataset"
	_ "github.com/rai-project/dldataset/vision"
	"github.com/rai-project/pipeline"
)

type getDataset struct {
	base
	dataset dldataset.Dataset
}

func NewGetDataset(dataset dldataset.Dataset) pipeline.Step {
	res := getDataset{
		dataset: dataset,
		base: base{
			info: "ListDataset",
		},
	}
	res.doer = res.do
	return res
}

func (p getDataset) do(ctx context.Context, in0 interface{}) interface{} {
	lst, err := p.dataset.List(ctx)
	if err != nil {
		return err
	}
	return lst
}

func (p getDataset) Close() error {
	return nil
}
