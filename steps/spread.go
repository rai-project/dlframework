package steps

import (
	"golang.org/x/net/context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/pipeline"
)

type spread struct {
	base
}

func NewSpread() pipeline.Step {
	step := spread{
		base: base{
			info:         "Spread",
			spreadOutput: true,
		},
	}
	step.doer = step.do
	return step
}

func (p spread) do(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{} {
	span, ctx := opentracing.StartSpanFromContext(ctx, p.Info())
	defer span.Finish()

	in, err := toSlice(in0)
	if err != nil {
		return errors.Errorf("expecting a slice input for Spread, but got %v", in0)
	}
	return in
}

func (p spread) Close() error {
	return nil
}
