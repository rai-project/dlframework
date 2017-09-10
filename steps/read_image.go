package steps

import (
	"io"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/image"
	"github.com/rai-project/pipeline"
)

type readImage struct {
	base
	options predict.PreprocessOptions
}

func NewReadImage(options predict.PreprocessOptions) pipeline.Step {
	res := readImage{
		base: base{
			info: "ReadImage",
		}}
	res.doer = res.do
	return res
}

func (p readImage) do(ctx context.Context, in0 interface{}) interface{} {
	in, ok := in0.(io.Reader)
	if !ok {
		return errors.Errorf("expecting a io.Reader for read image step, but got %v", in0)
	}

	width, height := 0, 0
	if len(p.options.Size) == 1 {
		width = p.options.Size[0]
		height = p.options.Size[0]
	}
	if len(p.options.Size) > 1 {
		width = p.options.Size[0]
		height = p.options.Size[1]
	}

	image, err := image.Read(in, image.Resized(width, height))
	if err != nil {
		return errors.Errorf("unable to read image")
	}

	return image
}
