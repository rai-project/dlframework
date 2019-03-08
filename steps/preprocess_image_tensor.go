package steps

import (
	"context"

	"github.com/k0kubun/pp"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
	"gorgonia.org/tensor"
)

type preprocessImageTensor struct {
	base
	options predictor.PreprocessOptions
}

func NewPreprocessImageTensor(options predictor.PreprocessOptions) pipeline.Step {
	res := preprocessImageTensor{
		base: base{
			info: "preprocessImageTensor",
		},
	}

	mean := []float32{0, 0, 0}
	if len(options.MeanImage) != 0 {
		mean = options.MeanImage
	}
	scale := float32(1.0)
	if options.Scale != 0 {
		scale = options.Scale
	}
	mode := types.RGBMode
	if options.ColorMode != 0 {
		mode = options.ColorMode
	}

	res.options = predictor.PreprocessOptions{
		Context:   options.Context,
		MeanImage: mean,
		Scale:     scale,
		ColorMode: mode,
		Layout:    options.Layout,
	}

	res.doer = res.do

	return res
}

func (p preprocessImageTensor) do(ctx context.Context, in0 interface{}, pipelineOptions *pipeline.Options) interface{} {
	if p.options.Context != nil {
		span, ctx0 := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, p.Info())
		ctx = ctx0
		defer span.Finish()
	}

	switch in := in0.(type) {
	case *types.RGBImage:
		if p.options.Layout == image.HWCLayout {
			return p.doRGBImageHWC(ctx, in)
		}
		panic("not implemented")
	case *types.BGRImage:
		panic("not implemented")
	}
	return errors.Errorf("expecting an RGB or BGR image for preprocess image step, but got %v", in0)
}

func (p preprocessImageTensor) doRGBImageHWC(ctx context.Context, in *types.RGBImage) interface{} {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	channels := 3
	scale := p.options.Scale
	// mode := p.options.ColorMode
	mean := p.options.MeanImage

	toFloat32 := func(x uint8) float32 {
		return float32(x)
	}
	tnsr, err := tensor.New(
		tensor.WithShape(height, width, channels),
		tensor.WithBacking(in.Pix),
		tensor.Of(tensor.Uint8),
	).Apply(interface{}(toFloat32))
	if err != nil {
		pp.Println(err)
		return errors.Wrapf(err, "unable to read the image")
	}

	pp.Println(tnsr.At(100, 100, 0))
	_ = scale
	_ = mean
	// if true {
	// 	tnsr.Sub(
	// 		tensor.New(
	// 			tensor.WithShape(3),
	// 			tensor.WithBacking(mean),
	// 			tensor.Of(tensor.Float32),
	// 		),
	// 		tensor.UseUnsafe(),
	// 	)
	// }
	// if false {
	// 	tnsr.DivScalar(
	// 		scale,
	// 		false,
	// 		tensor.UseUnsafe(),
	// 	)
	// }
	return tnsr
}

func (p preprocessImageTensor) Close() error {
	return nil
}
