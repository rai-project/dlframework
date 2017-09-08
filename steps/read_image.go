package steps

import (
	"io"

	"github.com/k0kubun/pp"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/rai-project/image"
	"github.com/rai-project/pipeline"
)

type readImage struct {
	base
}

func NewReadImage() pipeline.Step {
	res := readImage{}
	res.doer = res.do
	return res
}

func (p readImage) Info() string {
	return "ReadImage"
}

func (p readImage) do(ctx context.Context, in0 interface{}) interface{} {
	in, ok := in0.(io.Reader)
	if !ok {
		return errors.Errorf("expecting a io.Reader for read image step, but got %v", in0)
	}

	image, err := image.Read(ctx, in)
	if err != nil {
		pp.Printf("err = %v", err)
		return errors.Errorf("unable to read image")
	}

	return image
}
