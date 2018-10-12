package steps

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rai-project/dldataset"
	_ "github.com/rai-project/dldataset/vision"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
)

type getDataset struct {
	base
	dataset dldataset.Dataset
}

func NewGetDataset(dataset dldataset.Dataset) pipeline.Step {
	res := getDataset{
		dataset: dataset,
		base: base{
			info: "GetDataset",
		},
	}
	res.doer = res.do
	return res
}

func (p getDataset) do(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{} {
	if span, newCtx := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, p.Info()); span != nil {
		ctx = newCtx
		defer span.Finish()
	}

	in, ok := in0.(string)
	if !ok {
		return errors.Errorf("expecting a string for get dataset step, but got %v", in0)
	}

	lbl, err := p.dataset.Get(ctx, in)
	if err != nil {
		return err
	}

	return lbl
}

func (p getDataset) Close() error {
	return nil
}
