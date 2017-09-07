package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/pipeline"
)

type readImage struct {
	base
}

func NewReadImage() pipeline.Step {
	return readImage{}
}

func (p readImage) Info() string {
	return "ReadImage"
}

func (p readImage) do(ctx context.Context, in0 interface{}) interface{} {
	return nil
}

func (p readURL) Close() error {
	return nil
}
