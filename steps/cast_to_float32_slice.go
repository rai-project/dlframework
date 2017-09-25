package steps

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/rai-project/pipeline"
)

type castToFloat32Slice struct {
	base
}

func NewCastToFloat32Slice() pipeline.Step {
	p := castToFloat32Slice{
		base: base{
			info: "CastToFloat32Slice",
		},
	}
	p.doer = p.do
	return p
}

func (p castToFloat32Slice) do(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{} {
	span, ctx := opts.Tracer.StartSpanFromContext(ctx, p.Info())
	defer span.Finish()

	in, err := toSlice(in0)
	if err != nil {
		return errors.Errorf("expecting a slice input for CastToFloat32Slice, but got %v", in0)
	}
	res, err := toFloat32Slice(in)
	if err != nil {
		return err
	}
	return res
}

func (p castToFloat32Slice) Close() error {
	return nil
}
