package steps

import (
	"fmt"

	"github.com/facebookgo/stack"
	"github.com/fatih/color"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/uuid"
	"golang.org/x/net/context"
)

var (
	StackSize       = 4 << 10 // 4 KB
	DisableStackAll = false
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
	opts = append([]pipeline.Option{pipeline.Tracer(tracer)}, opts...)
	options := pipeline.NewOptions(opts...)
	go func() {
		defer close(out)
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = fmt.Errorf("%v", r)
				}
				stack := stack.Callers(3)
				log.WithError(err).WithField("step", p.Info()).Errorf("[%s] %v\n", color.RedString("PANIC RECOVER"), stack)
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
					id = a.GetID()
				} else {
					id = uuid.NewV4()
					// pp.Println("no id for %v @ step = %v", input, p.info)
				}

				res := p.doer(ctx, input, options)

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
