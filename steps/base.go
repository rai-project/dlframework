package steps

import (
	"github.com/rai-project/uuid"
	"golang.org/x/net/context"
)

type base struct {
	spreadOutput bool
	info         string
	doer         func(ctx context.Context, in0 interface{}) interface{}
}

func (p base) Info() string {
	return p.info
}

func (p base) Run(ctx context.Context, in <-chan interface{}, out chan interface{}) {
	go func() {
		defer close(out)
		defer func() {
			if r := recover(); r != nil {
				log.WithField("step", p.Info()).Error("recovered in %+v", r)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					out <- err
				}
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

				var id string

				org := input
				if a, ok := org.(IDer); ok {
					input = a.GetData()
					id = a.GetId()
				} else {
					id = uuid.NewV4()
					// pp.Println("no id for %v @ step = %v", input, p.info)
				}

				res := p.doer(ctx, input)

				if lst, ok := res.([]interface{}); ok && p.spreadOutput {
					// flatten sequence
					for _, e := range lst {
						out <- NewIDWrapper(id, e)
					}
					continue
				}

				out <- NewIDWrapper(id, res)

			}
		}
	}()
}

func (p base) Close() error {
	return nil
}
