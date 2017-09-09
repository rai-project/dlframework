package steps

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
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
	in, ok := in0.([]float32)
	if !ok {
		return errors.Errorf("expecting []float32 for predict image step, but got %v", in0)
	}

	features, err := p.predictor.Predict(ctx, in)
	if err != nil {
		return err
	}

	return features
}

func (p predictImage) Close() error {
	return nil
}
