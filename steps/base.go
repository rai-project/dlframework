package steps

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
)

type base struct{}

func (p base) do(ctx context.Context, in0 interface{}) interface{} {
	return errors.New("the base step is not implemented")
}

func (p base) Run(ctx context.Context, in <-chan interface{}, out chan interface{}) {
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				// if err := ctx.Err(); err != nil {
				// 	out <- err
				// }
				return
			case input, open := <-in:
				if !open {
					return
				}
				if err, ok := input.(error); ok {
					out <- err
					continue
				}

				org := input

				if a, ok := org.(IDer); ok {
					input = a.GetData()
				}

				res := p.do(ctx, input)

				if a, ok := org.(IDer); ok {
					res = NewIDWrapper(a.GetId(), res)
				}

				out <- res
			}
		}
	}()
}
