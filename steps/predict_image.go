package steps

import (
	"golang.org/x/net/context"

	opentracing "github.com/opentracing/opentracing-go"
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

func (p predictImage) do(ctx context.Context, in0 interface{}) interface{} {
	if span, newCtx := opentracing.StartSpanFromContext(ctx, p.Info()); span != nil {
		ctx = newCtx
		if framework, model, err := p.predictor.Info(); err == nil {
			span.SetTag("framework", framework.MustCanonicalName())
			span.SetTag("model", model.MustCanonicalName())
		}
		defer span.Finish()
	}

	in, ok := in0.([]float32)
	if !ok {
		return errors.Errorf("expecting []float32 for predict image step, but got %v", in0)
	}

	if p.predictor == nil {
		return errors.New("the predict image was created with a nil predictor")
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
