package steps

import (
	"context"
	"io"
	"io/ioutil"
	"strings"

	"github.com/h2non/filetype"
	"github.com/k0kubun/pp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/image"
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
		span, ctx := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, p.Info(), opentracing.Tags{
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
  
	croppedImg, err := cutter.Crop(in,  cutter.Config{
    Width: width*ratio,
    Height: height*ratio,
    Mode: p.options.CropMethod,
  })
	if err != nil {
		return errors.Errorf("unable to crop image")
	}

	return croppedImgs
}
