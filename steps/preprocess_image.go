package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/pipeline"
)

type preprocessImage struct {
	base
}

func NewPreprocessImage() pipeline.Step {
	return preprocessImage{}
}

func (p preprocessImage) Info() string {
	return "PreprocessImage"
}

func (p preprocessImage) do(ctx context.Context, in0 interface{}) interface{} {
	return nil
}

func (p preprocessImage) Close() error {
	return nil
}
