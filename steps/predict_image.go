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
	in, ok := in0.([]interface{})
	if !ok {
		return errors.Errorf("expecting []interface{} for predict image step, but got %v", in0)
	}

	var data [][]float32
	for _, e := range in {
		v, ok := e.([]float32)
		if !ok {
			return errors.Errorf("expecting []float32 for each image in predict image step, but got %v", e)
		}
		data = append(data, v)
	}

	if p.predictor == nil {
		return errors.New("the predict image was created with a nil predictor")
	}

	opts, err := p.predictor.GetPredictionOptions(ctx)
	if err != nil {
		return err
	}

	features, err := p.predictor.Predict(ctx, data, opts)
	if err != nil {
		return err
	}

	lst := make([]interface{}, len(data))
	for ii := 0; ii < len(in); ii++ {
		lst[ii] = features[ii]
	}

	return lst
}

func (p predictImage) Close() error {
	return nil
}
