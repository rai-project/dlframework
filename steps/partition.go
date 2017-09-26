package steps

import (
	"golang.org/x/net/context"

	"github.com/rai-project/pipeline"
)

type partition struct {
	base
	size int
}

func NewPartition(size int) pipeline.Step {
	res := partition{
		size: size,
	}
	return res
}

func (p partition) Info() string {
	return "Partition"
}

func (p partition) Run(ctx context.Context, in <-chan interface{}, out chan interface{}, opts ...pipeline.Option) {
	opts = append([]pipeline.Option{pipeline.Tracer(tracer)}, opts...)
	options := pipeline.NewOptions(opts...)
	_ = options
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
				panic("todo...")
				_ = input
				_ = open
			}
		}
	}()
}

func (p partition) Close() error {
	return nil
}
