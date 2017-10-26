package steps

import (
	"bytes"
	"io"
	"net/http"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
)

type readURL struct {
	base
}

func NewReadURL() pipeline.Step {
	res := readURL{
		base: base{
			info: "ReadURL",
		},
	}
	res.doer = res.do
	return res
}

func (p readURL) do(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{} {
	span, ctx := tracer.StartSpanFromContext(ctx, tracer.STEP_TRACE, p.Info())
	defer span.Finish()

	url := ""
	switch in := in0.(type) {
	case string:
		url = in
	case *dl.URLsRequest_URL:
		if in == nil {
			return errors.New("cannot read nil url")
		}
		url = in.GetData()
	case dl.URLsRequest_URL:
		url = in.GetData()
	default:
		return errors.Errorf("expecting a string for read url Step, but got %v", in0)
	}

	if span != nil {
		span.SetTag("url", url)
	}

	resp, err := http.Get(url)
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
