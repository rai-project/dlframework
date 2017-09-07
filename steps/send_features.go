package steps

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/pipeline"
	"golang.org/x/net/context"
)

type sendFeature struct {
	base
}

func NewSendFeature() (pipeline.Step, error) {
	return sendFeature{}, nil
}

func (p sendFeature) Info() string {
	return "SendFeature"
}

func (p sendFeature) do(ctx context.Context, in0 interface{}) interface{} {

	int16 := input

	if a, ok := org.(IDer); ok {
		in = a.GetData()
	}

	in, ok := in0.(dl.Features)
	if !ok {
		return errors.Errorf("expecting a string for read url Step, but got %v", in0)
	}
	resp, err := http.Get(in.GetUrl())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.Errorf("bad response code: %d", resp.StatusCode)
	}

	res := new(bytes.Buffer)
	_, err = io.Copy(res, resp.Body)
	if err != nil {
		return errors.Errorf("unable to copy body")
	}
	return res
}

func (p readURL) Close() error {
	return nil
}
