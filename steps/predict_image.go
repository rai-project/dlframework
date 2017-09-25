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
	res := predictImage{
		base: base{
			info: "PredictImage",
		},
	}
	res.predictor = predictor
	res.doer = res.do

	return res
}

func (p predictImage) do(ctx context.Context, in0 interface{}, pipelineOpts *pipeline.Options) interface{} {

	in, ok := in0.([]float32)
	if !ok {
		return errors.Errorf("expecting []float32 for predict image step, but got %v", in0)
	}

	if p.predictor == nil {
		return errors.New("the predict image was created with a nil predictor")
	}

	opts, err := p.predictor.GetPredictionOptions(ctx)
	if err != nil {
		return err
	}

	features, err := p.predictor.Predict(ctx, in, opts)
	if err != nil {
		return err
	}

	return features
}

func (p predictImage) Close() error {
	return nil
}
