package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/pipeline"
)

type asBGRImage struct {
	base
}

func NewAsBGRImage() pipeline.Step {
	res := asBGRImage{}
	res.doer = res.do
	return res
}

func (p asBGRImage) Info() string {
	return "AsBGRImage"
}

func (p asBGRImage) do(ctx context.Context, in0 interface{}) interface{} {
	return nil
}

func (p asBGRImage) Close() error {
	return nil
}
