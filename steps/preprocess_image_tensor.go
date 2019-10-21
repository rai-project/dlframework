package steps

import (
	"context"

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
			info: "preprocess_image_tensor_step",
		},
		options: options,
	}
	res.doer = res.do
	return res
}

func (p preprocessImageTensor) do(ctx context.Context, in0 interface{}, pipelineOptions *pipeline.Options) interface{} {
	span, _ := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, p.Info())
	defer span.Finish()

	switch in := in0.(type) {
	case *types.RGBImage:
		if p.options.Layout == image.HWCLayout {
			return p.doRGBImageHWC(in)
		}
		panic("not implemented")
	case *types.BGRImage:
		panic("not implemented")
	}
	return errors.Errorf("expecting an RGB or BGR image for preprocess image step, but got %v", in0)
}

func (p preprocessImageTensor) doRGBImageHWC(in *types.RGBImage) interface{} {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	channels := 3
	// scale := p.options.Scale
	// mode := p.options.ColorMode
	// mean := p.options.MeanImage

	tnsrBytes := tensor.New(
		tensor.WithShape(height, width, channels),
		tensor.WithBacking(in.Pix),
	)

	// toFloat32 := func(x uint8) float32 {
	// 	return float32(x)
	// }
	// tnsrFloats, err := tnsrBytes.Apply(interface{}(toFloat32))
	// if err != nil {
	// 	return errors.Wrapf(err, "unable to read the image")
	// }

	backing := tnsrBytes.Data().([]byte)
	backing2 := make([]float32, len(backing))
	for i := range backing {
		backing2[i] = float32(backing[i])
	}
	tnsrFloats := tensor.New(tensor.WithShape(tnsrBytes.Shape()...), tensor.WithBacking(backing2))

	// tnsrFloats.Sub(
	// 	tensor.New(
	// 		tensor.WithShape(3),
	// 		tensor.WithBacking(mean),
	// 		tensor.Of(tensor.Float32),
	// 	),
	// 	tensor.UseUnsafe(),
	// )

	// if false {
	// 	tnsrFloats.DivScalar(
	// 		scale,
	// 		false,
	// 		tensor.UseUnsafe(),
	// 	)
	// }
	return tnsrFloats
}

func (p preprocessImageTensor) Close() error {
	return nil
}
