package steps

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/h2non/filetype"
	"github.com/k0kubun/pp"

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
