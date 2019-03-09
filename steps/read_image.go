package steps

import (
	"io"

	"context"

	"github.com/pkg/errors"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
)

type readImage struct {
	base
	options predictor.PreprocessOptions
}

func NewReadImage(options predictor.PreprocessOptions) pipeline.Step {
	res := readImage{
		base: base{
			info: "ReadImage",
		},
		options: options,
	}
	res.doer = res.do
	return res
}

func (p readImage) do(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{} {
	// no need to trace here, since resize and read already perform tracing

	if p.options.Context == nil {
		ctx = nil
	}
	readOptions := []image.Option{
		image.Context(ctx),
	}
	dims := p.options.Dims
	if dims != nil {
		height, width := dims[1], dims[2]
		readOptions = append(readOptions, image.Resized(height, width))
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

	image, err := image.Read(in, readOptions...)
	if err != nil {
		return errors.Errorf("unable to read image")
	}
	return image

}
