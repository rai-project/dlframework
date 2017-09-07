package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/pipeline"
)

type asRGBImage struct {
	base
}

func NewAsRGBImage() pipeline.Step {
	return asRGBImage{}
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
