package steps

import (
	"io"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
)

type readImage struct {
	base
	width   int
	height  int
	options predict.PreprocessOptions
}

func NewReadImage(options predict.PreprocessOptions) pipeline.Step {

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

func (p readImage) do(ctx context.Context, in0 interface{}) interface{} {
	var in io.Reader

	switch e := in0.(type) {
	case io.Reader:
		in = e
	case dldataset.LabeledData:
		data, err := e.Data()
		if err != nil {
			return err
		}
		switch d := data.(type) {
		case *types.RGBImage:
			img, err := image.Resize(d, p.width, p.height)
			if err != nil {
				return err
			}
			return img
		case *types.BGRImage:
			img, err := image.Resize(d, p.width, p.height)
			if err != nil {
				return err
			}
			return img
		case io.Reader:
			in = d
		default:
			return errors.Errorf("expecting a io.Reader or image for read image step, but got %v", in0)

		}
	default:
		return errors.Errorf("expecting a io.Reader or dataset element for read image step, but got %v", in0)
	}

	image, err := image.Read(in, image.Resized(p.width, p.height))
	if err != nil {
		return errors.Errorf("unable to read image")
	}

	return image
}
