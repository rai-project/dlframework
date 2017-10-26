package steps

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
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
	span, ctx := tracer.StartSpanFromContext(ctx, tracer.STEP_TRACE, p.Info())
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
