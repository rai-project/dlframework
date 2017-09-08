package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/dldataset"
	_ "github.com/rai-project/dldataset/vision"
	"github.com/rai-project/pipeline"
)

type listDataset struct {
	base
	category string
	name     string
	dataset  dldataset.Dataset
}

func NewListDataset(category, name string) pipeline.Step {
	dataset, err := dldataset.Get(category, name)
	if err != nil {
		panic(err)
	}
	res := listDataset{
		category: category,
		name:     name,
		dataset:  dataset,
	}
	res.doer = res.do
	return res
}

func (p listDataset) Info() string {
	return "ListDataset"
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
