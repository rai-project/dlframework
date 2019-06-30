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

type readImage struct {
	base
	options predictor.PreprocessOptions
}

func NewReadImage(options predictor.PreprocessOptions) pipeline.Step {
	res := readImage{
		base: base{
			info: "read_image_step",
		},
		options: options,
	}
	res.doer = res.do
	return res
}

func (p readImage) do(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{} {
	readOptions := []image.Option{}
	if opentracing.SpanFromContext(ctx) != nil {
		span, ctx := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, p.Info(), opentracing.Tags{
			"trace_source": "steps",
			"step_name":    "read_image",
		})
		defer span.Finish()
		readOptions = []image.Option{
			image.Context(ctx),
		}
	}

	readOptions = append(readOptions, image.Mode(p.options.ColorMode))

	dims := p.options.Dims
	if dims != nil && len(dims) > 2 {
		height, width := dims[1], dims[2]
		readOptions = append(readOptions, image.Resized(height, width))
	}
	if p.options.MaxDimension != nil {
		readOptions = append(readOptions, image.MaxDimension(*p.options.MaxDimension))
	}
	if p.options.KeepAspectRatio != nil {
		readOptions = append(readOptions, image.KeepAspectRatio(*p.options.KeepAspectRatio))
	}

	var in io.Reader
	switch in0 := in0.(type) {
	case io.Reader:
		in = in0
	case dldataset.LabeledData:
		data, err := in0.Data()
		if err != nil {
			return err
		}
		switch data := data.(type) {
		case *types.RGBImage:
			img, err := image.Resize(data, readOptions...)
			if err != nil {
				return err
			}
			return img
		case *types.BGRImage:
			img, err := image.Resize(data, readOptions...)
			if err != nil {
				return err
			}
			return img
		case io.Reader:
			in = data
		default:
			return errors.Errorf("expecting a io.Reader or image for read image step, but got %v", in0)

		}
	default:
		return errors.Errorf("expecting a io.Reader or dataset element for read image step, but got %v", in0)
	}

	elementType := strings.ToLower(p.options.ElementType)

	if elementType == "raw_image" {
		reader, ok := in.(io.Reader)
		if !ok {
			return errors.Errorf("expecting an io.Reader data type for %v step, but got %v", p.Info(), pp.Sprint(in0))
		}
		buf, err := ioutil.ReadAll(reader)
		if err != nil {
			return errors.Wrapf(err, "failed to read data for %v step", p.Info())
		}
		if !filetype.IsImage(buf[:261]) {
			return errors.Errorf("expecting a raw image for %v step, but got %v", p.Info(), pp.Sprint(in0))
		}
		return buf
	}

	image, err := image.Read(in, readOptions...)
	if err != nil {
		return errors.Errorf("unable to read image")
	}

	return image
}
