package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/pipeline"
)

type preprocessImage struct {
	base
	options predict.PreprocessOptions
}

func NewPreprocessImage() pipeline.Step {
	res := preprocessImage{}
	res.doer = res.do
	return res
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
