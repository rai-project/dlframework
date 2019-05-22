package steps

import (
	"context"

	"github.com/oliamb/cutter"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
)

type cropImage struct {
	base
	options predictor.PreprocessOptions
}

func NewCropImage(options predictor.PreprocessOptions) pipeline.Step {
	res := cropImage{
		base: base{
			info: "CropImageStep",
		},
		options: options,
	}
	res.doer = res.do
	return res
}

func (p cropImage) do(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{} {
	if opentracing.SpanFromContext(ctx) != nil {
		span, _ := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, p.Info(), opentracing.Tags{
			"trace_source": "steps",
			"step_name":    "crop_image",
		})
		defer span.Finish()
	}

	in, ok := in0.(types.Image)
	if !ok {
		return errors.Errorf("expecting a io.Reader or dataset element for read image step, but got %v", in0)
	}

	ratio := p.options.CropRatio
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()

	croppedImg, err := cutter.Crop(in, cutter.Config{
		Width:  int(float32(width) * ratio),
		Height: int(float32(height) * ratio),
		Mode:   p.options.CropMethod,
	})
	if err != nil {
		return errors.Errorf("unable to crop image")
	}

	return croppedImg
}
