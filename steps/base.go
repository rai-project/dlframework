package steps

import (
	"golang.org/x/net/context"
)

type base struct {
	spreadOutput bool
	doer         func(ctx context.Context, in0 interface{}) interface{}
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
				// pp.Printf("input = %v, open = %v\n", input, open)
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

				res := p.doer(ctx, input)
				if lst, ok := res.([]interface{}); ok && p.spreadOutput {
					// flatten sequence
					for _, e := range lst {
						if a, ok := org.(IDer); ok {
							e = NewIDWrapper(a.GetId(), e)
						}
						out <- e
					}
				} else {
					if a, ok := org.(IDer); ok {
						res = NewIDWrapper(a.GetId(), res)
					}
					out <- res
				}

			}
		}
	}()
}

func (p base) Close() error {
	return nil
}
