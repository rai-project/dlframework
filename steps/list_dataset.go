package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/dldataset"
	_ "github.com/rai-project/dldataset/vision"
	"github.com/rai-project/pipeline"
)

type listDataset struct {
	base
	dataset dldataset.Dataset
}

func NewListDataset(dataset dldataset.Dataset) pipeline.Step {
	res := listDataset{
		dataset: dataset,
		base: base{
			info: "ListDataset",
		},
	}
	res.doer = res.do
	return res
}

func (p listDataset) do(ctx context.Context, in0 interface{}) interface{} {
	lst, err := p.dataset.List(ctx)
	if err != nil {
		return err
	}
	return lst
}

func (p listDataset) Close() error {
	return nil
}
