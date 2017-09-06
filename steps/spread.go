package steps

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/rai-project/pipeline"
)

type spread struct {
	base
}

func NewSpread(f func(in interface{}) interface{}) (pipeline.Step, error) {
	return spread{
		base: base{
			spreadOutput: true,
		},
	}, nil
}

func (p spread) Info() string {
	return "Spread"
}

func (p spread) do(ctx context.Context, in0 interface{}) interface{} {
	in, err := toSlice(in0)
	if err != nil {
		return errors.Errorf("expecting a slice input for Spread, but got %v", in0)
	}
	return in
}

func (p spread) Close() error {
	return nil
}
