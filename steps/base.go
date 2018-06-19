package steps

import (
	"context"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/uuid"
)

type base struct {
	spreadOutput bool
	info         string
	doer         func(ctx context.Context, in0 interface{}, opts *pipeline.Options) interface{}
}

func (p base) Info() string {
	return p.info
}

func (p base) Run(ctx context.Context, in <-chan interface{}, out chan interface{}, opts ...pipeline.Option) {
	options := pipeline.NewOptions(opts...)
	go func() {
		defer close(out)
		defer onPanic(p.Info())
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

				var id, data interface{}

				org := input
				switch a := org.(type) {
				case IDer:
					data = a.GetData()
					id = a.GetID()
				case []interface{}:
					tempID := make([]string, len(a))
					tempData := make([]interface{}, len(a))
					for ii, e := range a {
						switch e := e.(type) {
						case IDer:
							tempID[ii] = e.GetID()
							tempData[ii] = e.GetData()
						default:
							tempID[ii] = uuid.NewV4()
							tempData[ii] = e
						}
					}
					id = tempID
					data = tempData
				default:
					id = uuid.NewV4()
					data = input
					// pp.Println("no id for %v @ step = %v", input, p.info)
				}

				res := p.doer(ctx, data, options)

				if err, ok := res.(error); ok {
					out <- err
					continue
				}

				if lst, ok := res.([]interface{}); ok && p.spreadOutput {
					// flatten sequence
					for _, e := range lst {
						out <- NewIDWrapper(id.(string), e)
					}
					continue
				}

				switch ids := id.(type) {
				case []string:
					lst, ok := res.([]interface{})
					if !ok {
						log.Errorf("expecting a []inteface{} for result but got %v", res)
						return
					}

					for ii, e := range lst {
						out <- NewIDWrapper(ids[ii], e)
					}
				default:
					out <- NewIDWrapper(id.(string), res)
				}

			}
		}
	}()
}

func (p base) Close() error {
	return nil
}
