package steps

import (
	"context"

	cupti "github.com/rai-project/go-cupti"
	"github.com/rai-project/tracer"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/options"
	"github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/pipeline"
)

type predictObjectDetection struct {
	base
	predictor predictor.Predictor
}

func NewPredictObjectDetection(predictor predictor.Predictor) pipeline.Step {
	res := predictObjectDetection{
		base: base{
			info: "PredictObjectDetection",
		},
	}
	res.predictor = predictor
	res.doer = res.do

	return res
}

func (p predictObjectDetection) do(ctx context.Context, in0 interface{}, pipelineOpts *pipeline.Options) interface{} {
	data, ok := in0.([]interface{})
	if !ok {
		return errors.Errorf("expecting []interface{} for predict image step, but got %v", in0)
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
	lst := make([]interface{}, len(data))
	for ii := 0; ii < len(data); ii++ {
		lst[ii] = features[ii]
	}

	return lst
}

func (p predictObjectDetection) Close() error {
	return nil
}
