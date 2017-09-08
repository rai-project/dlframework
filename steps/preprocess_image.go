package steps

import (
	"golang.org/x/net/context"

	"github.com/anthonynsimon/bild/parallel"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
)

type preprocessImage struct {
	base
	options predict.PreprocessOptions
}

func NewPreprocessImage(options predict.PreprocessOptions) pipeline.Step {
	res := preprocessImage{}

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
	res.options = predict.PreprocessOptions{
		MeanImage: mean,
		Scale:     scale,
		ColorMode: mode,
	}

	res.doer = res.do

	return res
}

func (p preprocessImage) Info() string {
	return "PreprocessImage"
}

func (p preprocessImage) do(ctx context.Context, in0 interface{}) interface{} {
	in, ok := in0.(*types.RGBImage)
	if !ok {
		return errors.Errorf("expecting a predict.PreprocessOptions for preprocess image step, but got %v", in0)
	}

	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage

	out := make([]float32, 3*height*width)
	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				offset := y*in.Stride + x*3
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				if mode == types.RGBMode {
					out[y*width+x] = (float32(r) - mean[0]) / scale
					out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
					out[2*width*height+y*width+x] = (float32(b) - mean[2]) / scale
				} else if mode == types.BGRMode {
					out[y*width+x] = (float32(b) - mean[2]) / scale
					out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
					out[2*width*height+y*width+x] = (float32(r) - mean[0]) / scale
				} else {
					panic("invalid mode in preprocess image step")
				}
			}
		}
	})

	return out
}

func (p preprocessImage) Close() error {
	return nil
}
