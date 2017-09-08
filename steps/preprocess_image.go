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

func NewPreprocessImage() pipeline.Step {
	res := preprocessImage{}
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

	height := p.options.Size[0]
	width := p.options.Size[1]
	scale := p.options.Scale
	mode := p.options.ColorSpace
	mean := p.options.MeanImage

	out := make([]uint8, 3*height*width)
	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				rgb := in.RGBAt(x, y)
				if mode == types.RGBMode {
					out[y*width+x] = uint8((float32(rgb.R) - mean[0]) / scale)
					out[width*height+y*width+x] = uint8((float32(rgb.G) - mean[1]) / scale)
					out[2*width*height+y*width+x] = uint8((float32(rgb.B) - mean[2]) / scale)
				} else if mode == types.BGRMode {
					out[y*width+x] = uint8((float32(rgb.B) - mean[2]) / scale)
					out[width*height+y*width+x] = uint8((float32(rgb.G) - mean[1]) / scale)
					out[2*width*height+y*width+x] = uint8((float32(rgb.R) - mean[0]) / scale)
				} else {
					// TODO
				}
			}
		}
	})

	return out
}

func (p preprocessImage) Close() error {
	return nil
}
