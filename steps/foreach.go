package steps

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/rai-project/pipeline"
)

type forEach struct {
	base
	f func(in interface{}) interface{}
}

func NewForEach(ctx context.Context, f func(in interface{}) interface{}) (pipeline.Step, error) {
	return forEach{
		f: f,
	}, nil
}

func (p forEach) New(ctx context.Context) (pipeline.Step, error) {
	return p, nil
}

func (p forEach) Info() string {
	return "ForEach"
}

func (p forEach) do(ctx context.Context, in0 interface{}) interface{} {
	in, err := toSlice(in0)
	if err != nil {
		return errors.Errorf("expecting a slice input for CastToFloat32Slice, but got %v", in0)
	}
	res := make([]interface{}, len(in))
	for ii, e := range in {
		res[ii] = p.f(e)
	}
	return res
}

func (p forEach) Close() error {
	return nil
}
