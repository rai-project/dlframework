package steps

import (
	"context"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/options"
	"github.com/rai-project/dlframework/framework/predictor"
	cupti "github.com/rai-project/go-cupti"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
)

var DefaultModelElementType string = "float32"

type predict struct {
	base
	predictor predictor.Predictor
}

func NewPredict(predictor predictor.Predictor) pipeline.Step {
	res := predict{
		base: base{
			info: "Predict",
		},
	}
	res.predictor = predictor
	res.doer = res.do

	return res
}

func (p predict) do(ctx context.Context, in0 interface{}, pipelineOpts *pipeline.Options) interface{} {
	iData, ok := in0.([]interface{})
	if !ok {
		return errors.Errorf("expecting []interface{} for predict image step, but got %v", in0)
	}

	data, err := p.castToElementType(iData)
	if err != nil {
		return err
	}

	if p.predictor == nil {
		return errors.New("the predict image was created with a nil predictor")
	}

	opts, err := p.predictor.GetPredictionOptions(ctx)
	if err != nil {
		return err
	}

	framework, model, err := p.predictor.Info()
	if err != nil {
		return err
	}

	span, ctx := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, p.Info(), opentracing.Tags{
		"model_name":        model.GetName(),
		"model_version":     model.GetVersion(),
		"framework_name":    framework.GetName(),
		"framework_version": framework.GetVersion(),
		"batch_size":        opts.BatchSize(),
		"feature_limit":     opts.FeatureLimit(),
		"device":            opts.Devices().String(),
		"trace_level":       opts.TraceLevel().String(),
		"uses_gpu":          opts.UsesGPU(),
	})
	defer span.Finish()

	var cu *cupti.CUPTI

	if opts.UsesGPU() && opts.TraceLevel() >= tracer.HARDWARE_TRACE {
		cu, err = cupti.New(cupti.Context(ctx))
	}

	err = p.predictor.Predict(ctx, data, options.WithOptions(opts))
	if err != nil {
		if cu != nil {
			cu.Wait()
			cu.Close()
		}
		return err
	}

	if cu != nil {
		cu.Wait()
		cu.Close()
	}

	features, err := p.predictor.ReadPredictedFeatures(ctx)
	lst := make([]interface{}, len(iData))
	for ii := 0; ii < len(iData); ii++ {
		lst[ii] = features[ii]
	}

	return lst
}

func (p predict) castToElementType(inputs []interface{}) (interface{}, error) {
	_, model, _ := p.predictor.Info()

	switch t := strings.ToLower(model.GetElementType(DefaultModelElementType)); t {
	case "uint8":
		data := make([][]uint8, len(inputs))
		for ii, input := range inputs {
			r, err := toUInt8Slice(input)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to cast to uint8 slice in %v step", p.info)
			}
			data[ii] = r
		}
		return data, nil
	case "float32":
		data := make([][]float32, len(inputs))
		for ii, input := range inputs {
			r, err := toFloat32Slice(input)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to cast to float32 slice in %v step", p.info)
			}
			data[ii] = r
		}
		return data, nil
	default:
		return nil, errors.Errorf("unsupported element type %v", t)
	}
}

func (p predict) Close() error {
	return nil
}
