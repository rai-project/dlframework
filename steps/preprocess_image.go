package steps

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
)

type preprocessImage struct {
	base
	options predictor.PreprocessOptions
}

func NewPreprocessImage(options predictor.PreprocessOptions) pipeline.Step {
	res := preprocessImage{
		base: base{
			info: "PreprocessImage",
		},
		options: options,
	}
	res.doer = res.do
	return res
}

func (p preprocessImage) do(ctx context.Context, in0 interface{}, pipelineOptions *pipeline.Options) interface{} {
	if p.options.Context != nil {
		span, ctx0 := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, p.Info())
		ctx = ctx0
		defer span.Finish()
	}

	var out []float32
	switch in := in0.(type) {
	case *types.RGBImage:
		return in.Pix
		if p.options.Layout == image.CHWLayout {
			out = p.doRGBImageCHW(ctx, in)
		}
		out = p.doRGBImageHWC(ctx, in)
	case *types.BGRImage:
		return in.Pix
		if p.options.Layout == image.CHWLayout {
			out = p.doBGRImageCHW(ctx, in)
		}
		out = p.doBGRImageHWC(ctx, in)
	default:
		return errors.Errorf("expecting an RGB or BGR image for preprocess image step, but got %v", in0)
	}

	elementType := strings.ToLower(p.options.ElementType)
	switch elementType {
	case "float32":
		return out
	case "uint8":
		out0 := make([]uint8, len(out))
		for ii, _ := range out {
			out0[ii] = uint8(out[ii])
		}
		return out0
	}

	return errors.Errorf("unsupported element type %v", elementType)
}

func (p preprocessImage) doRGBImageCHW(ctx context.Context, in *types.RGBImage) []float32 {
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()

	out := make([]float32, 3*height*width)

	if mode == types.RGBMode {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				offset := y*in.Stride + x*3
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				out[y*width+x] = (float32(r) - mean[0]) / scale
				out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
				out[2*width*height+y*width+x] = (float32(b) - mean[2]) / scale
			}
		}
	} else if mode == types.BGRMode {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				offset := y*in.Stride + x*3
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				out[y*width+x] = (float32(b) - mean[2]) / scale
				out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
				out[2*width*height+y*width+x] = (float32(r) - mean[0]) / scale
			}
		}
	} else {
		panic("invalid mode in preprocess image step")
	}

	return out
}

func (p preprocessImage) doRGBImageHWC(ctx context.Context, in *types.RGBImage) []float32 {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage

	out := make([]float32, 3*height*width)

	if mode == types.RGBMode {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				offset := y*in.Stride + x*3
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				out[offset+0] = (float32(r) - mean[0]) / scale
				out[offset+1] = (float32(g) - mean[1]) / scale
				out[offset+2] = (float32(b) - mean[2]) / scale
			}
		}
	} else if mode == types.BGRMode {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				offset := y*in.Stride + x*3
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				out[offset+0] = (float32(b) - mean[2]) / scale
				out[offset+1] = (float32(g) - mean[1]) / scale
				out[offset+2] = (float32(r) - mean[0]) / scale
			}
		}
	} else {
		panic("invalid mode in preprocess image step")
	}

	return out
}

func (p preprocessImage) doBGRImageCHW(ctx context.Context, in *types.BGRImage) []float32 {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage

	out := make([]float32, 3*height*width)
	if mode == types.RGBMode {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				offset := y*in.Stride + x*3
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				out[y*width+x] = (float32(b) - mean[0]) / scale
				out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
				out[2*width*height+y*width+x] = (float32(r) - mean[2]) / scale
			}
		}
	} else if mode == types.BGRMode {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				offset := y*in.Stride + x*3
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				out[y*width+x] = (float32(r) - mean[2]) / scale
				out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
				out[2*width*height+y*width+x] = (float32(b) - mean[0]) / scale
			}
		}
	} else {
		panic("invalid mode in preprocess image step")
	}

	return out
}

func (p preprocessImage) doBGRImageHWC(ctx context.Context, in *types.BGRImage) []float32 {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage

	out := make([]float32, 3*height*width)

	if mode == types.RGBMode {
		offset := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				out[offset+0] = (float32(b) - mean[0]) / scale
				out[offset+1] = (float32(g) - mean[1]) / scale
				out[offset+2] = (float32(r) - mean[2]) / scale
				offset = offset + 3
			}
		}
	} else if mode == types.BGRMode {
		offset := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				rgb := in.Pix[offset : offset+3]
				r, g, b := rgb[0], rgb[1], rgb[2]
				out[offset+0] = (float32(r) - mean[2]) / scale
				out[offset+1] = (float32(g) - mean[1]) / scale
				out[offset+2] = (float32(b) - mean[0]) / scale
				offset = offset + 3
			}
		}
	} else {
		panic("invalid mode in preprocess image step")
	}

	return out
}

func (p preprocessImage) Close() error {
	return nil
}
