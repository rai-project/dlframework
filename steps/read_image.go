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
	width   int
	height  int
	options predictor.PreprocessOptions
}

func NewReadImage(options predictor.PreprocessOptions) pipeline.Step {
	width, height := 0, 0
	if len(options.Size) == 1 {
		width = options.Size[0]
		height = options.Size[0]
	}
	if len(options.Size) > 1 {
		width = options.Size[0]
		height = options.Size[1]
	}

	res := readImage{
		width:   width,
		height:  height,
		options: options,
		base: base{
			info: "ReadImage",
		},
	}
	res.doer = res.do
	return res
}

func (p readImage) do(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{} {
	// no need to trace here, since resize and read already perform tracing

	var in io.Reader

	if p.options.Context == nil {
		ctx = nil
	}

	readOptions := []image.Option{
		image.Context(ctx),
	}

	if p.width != -1 || p.height != -1 {
		readOptions = append(readOptions, image.Resized(p.width, p.height))
	}

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
