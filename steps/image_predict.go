package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/pipeline"
)

type predictImage struct {
	base
	predictor predict.Predictor
}

func NewPredictImage(predictor predict.Predictor) pipeline.Step {
	res := predictImage{}
	res.predictor = predictor
	res.doer = res.do

	return res
}

func (p predictImage) Info() string {
	return "predictImage"
}

func (p predictImage) do(ctx context.Context, in0 interface{}) interface{} {
	return nil
}

func (p predictImage) Close() error {
	return nil
}
