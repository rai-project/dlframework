package steps

import (
	"golang.org/x/net/context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/pipeline"
)

type forEach struct {
	base
	f func(in interface{}) interface{}
}

func NewForEach(f func(in interface{}) interface{}) pipeline.Step {
	res := forEach{
		f: f,
	}
	res.doer = res.do
	return res
}

func (p forEach) Info() string {
	return "ForEach"
}

func (p forEach) do(ctx context.Context, in0 interface{}) interface{} {
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

func (p forEach) Close() error {
	return nil
}
