package steps

import (
	"golang.org/x/net/context"

	"github.com/anthonynsimon/bild/parallel"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
)

type preprocessImage struct {
	base
	options predict.PreprocessOptions
}

func NewPreprocessImage(options predict.PreprocessOptions) pipeline.Step {
	res := preprocessImage{
		base: base{
			info: "PreprocessImage",
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

	res.options = predict.PreprocessOptions{
		MeanImage: mean,
		Scale:     scale,
		ColorMode: mode,
		Layout: options.Layout,
	}

	res.doer = res.do

	return res
}

func (p preprocessImage) do(ctx context.Context, in0 interface{}, pipelineOptions *pipeline.Options) interface{} {
	span, ctx := pipelineOptions.Tracer.StartSpanFromContext(ctx, p.Info())
	defer span.Finish()

	switch in := in0.(type) {
	case *types.RGBImage:
		if p.options.Layout == image.CHWLayout {
			return p.doRGBImageCHW(ctx, in)
		}
		return p.doRGBImageHWC(ctx, in)
	case *types.BGRImage:
		if p.options.Layout == image.CHWLayout {
			return p.doBGRImageCHW(ctx, in)
		}
		return p.doBGRImageHWC(ctx, in)
	}
	return errors.Errorf("expecting an RGB or BGR image for preprocess image step, but got %v", in0)
}

func (p preprocessImage) doRGBImageCHW(ctx context.Context, in *types.RGBImage) interface{} {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage

	out := make([]float32, 3*height*width)

	if mode == types.RGBMode {
		parallel.Line(height, func(start, end int) {
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					offset := y*in.Stride + x*3
					rgb := in.Pix[offset : offset+3]
					r, g, b := rgb[0], rgb[1], rgb[2]
					out[y*width+x] = (float32(r) - mean[0]) / scale
					out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
					out[2*width*height+y*width+x] = (float32(b) - mean[2]) / scale
				}
			}
		})
	} else if mode == types.BGRMode {
		parallel.Line(height, func(start, end int) {
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					offset := y*in.Stride + x*3
					rgb := in.Pix[offset : offset+3]
					r, g, b := rgb[0], rgb[1], rgb[2]
					out[y*width+x] = (float32(b) - mean[2]) / scale
					out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
					out[2*width*height+y*width+x] = (float32(r) - mean[0]) / scale
				}
			}
		})
	} else {
		panic("invalid mode in preprocess image step")
	}

	return out
}

func (p preprocessImage) doRGBImageHWC(ctx context.Context, in *types.RGBImage) interface{} {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage

	out := make([]float32, 3*height*width)

	if mode == types.RGBMode {
		parallel.Line(height, func(start, end int) {
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					offset := y * in.Stride + x*3
					rgb := in.Pix[offset : offset+3]
					r, g, b := rgb[0], rgb[1], rgb[2]
					out[offset+0] = (float32(r) - mean[0]) / scale
					out[offset+1] = (float32(g) - mean[1]) / scale
					out[offset+2] = (float32(b) - mean[2]) / scale
				}
			}
		})
	} else if mode == types.BGRMode {
		parallel.Line(height, func(start, end int) {
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					offset := y * in.Stride + x*3
					rgb := in.Pix[offset : offset+3]
					r, g, b := rgb[0], rgb[1], rgb[2]
					out[offset+0] = (float32(b) - mean[2]) / scale
					out[offset+1] = (float32(g) - mean[1]) / scale
					out[offset+2] = (float32(r) - mean[0]) / scale
				}
			}
		})
	} else {
		panic("invalid mode in preprocess image step")
	}

	return out
}

func (p preprocessImage) doBGRImageCHW(ctx context.Context, in *types.BGRImage) interface{} {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage

	out := make([]float32, 3*height*width)
	if mode == types.RGBMode {
		parallel.Line(height, func(start, end int) {
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					offset := y*in.Stride + x*3
					rgb := in.Pix[offset : offset+3]
					r, g, b := rgb[0], rgb[1], rgb[2]
					out[y*width+x] = (float32(b) - mean[0]) / scale
					out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
					out[2*width*height+y*width+x] = (float32(r) - mean[2]) / scale
				}
			}
		})
	} else if mode == types.BGRMode {
		parallel.Line(height, func(start, end int) {
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					offset := y*in.Stride + x*3
					rgb := in.Pix[offset : offset+3]
					r, g, b := rgb[0], rgb[1], rgb[2]
					out[y*width+x] = (float32(r) - mean[2]) / scale
					out[width*height+y*width+x] = (float32(g) - mean[1]) / scale
					out[2*width*height+y*width+x] = (float32(b) - mean[0]) / scale
				}
			}
		})
	} else {
		panic("invalid mode in preprocess image step")
	}

	return out
}

func (p preprocessImage) doBGRImageHWC(ctx context.Context, in *types.BGRImage) interface{} {
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()
	scale := p.options.Scale
	mode := p.options.ColorMode
	mean := p.options.MeanImage

	out := make([]float32, 3*height*width)

	if mode == types.RGBMode {
		parallel.Line(height, func(start, end int) {
			offset := start * in.Stride
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					rgb := in.Pix[offset : offset+3]
					r, g, b := rgb[0], rgb[1], rgb[2]
					out[offset+0] = (float32(b) - mean[0]) / scale
					out[offset+1] = (float32(g) - mean[1]) / scale
					out[offset+2] = (float32(r) - mean[2]) / scale
					offset = offset + 3
				}
			}
		})
	} else if mode == types.BGRMode {
		parallel.Line(height, func(start, end int) {
			offset := start * in.Stride
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					rgb := in.Pix[offset : offset+3]
					r, g, b := rgb[0], rgb[1], rgb[2]
					out[offset+0] = (float32(r) - mean[2]) / scale
					out[offset+1] = (float32(g) - mean[1]) / scale
					out[offset+2] = (float32(b) - mean[0]) / scale
					offset = offset + 3
				}
			}
		})
	} else {
		panic("invalid mode in preprocess image step")
	}

	return out
}

func (p preprocessImage) Close() error {
	return nil
}
