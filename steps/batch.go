package steps

import (
	"golang.org/x/net/context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/pipeline"
)

type batch struct {
	base
	f     func(in interface{}) interface{}
	count int
}

func NewBatch(f func(in interface{}) interface{}, cnt int) pipeline.Step {
	res := batch{
		f:     f,
		count: cnt,
	}
	res.doer = res.do
	return res
}

func (p batch) Info() string {
	return "Batch"
}

func (p batch) do(ctx context.Context, in0 interface{}) interface{} {
	if span, newCtx := opentracing.StartSpanFromContext(ctx, p.Info()); span != nil {
		ctx = newCtx
		defer span.Finish()
	}

	in, err := toSlice(in0)
	if err != nil {
		return errors.Errorf("expecting a slice input for Spread, but got %v", in0)
	}
	res := make([]interface{}, len(in))
	for ii, e := range in {
		res[ii] = p.f(e)
	}
	return res
}

func (p batch) Close() error {
	return nil
}
