package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/pipeline"
)

type asRGBImage struct {
	base
}

func NewAsRGBImage() pipeline.Step {
	res := asRGBImage{}
	res.doer = res.do
	return res
}

func (p asRGBImage) Info() string {
	return "AsRGBImage"
}

func (p asRGBImage) do(ctx context.Context, in0 interface{}) interface{} {
	return nil
}

func (p asRGBImage) Close() error {
	return nil
}
