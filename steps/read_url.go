package steps

import (
	"bytes"
	"io"
	"net/http"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/pipeline"
)

type readURL struct {
	base
}

func NewReadURL() (pipeline.Step, error) {
	return readURL{}, nil
}

func (p readURL) Info() string {
	return "ReadURL"
}

func (p readURL) do(ctx context.Context, in0 interface{}) interface{} {
	in, ok := in0.(dl.URLsRequest_URL)
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