package options

import (
	dl "github.com/rai-project/dlframework"
	context "golang.org/x/net/context"
)

type Options struct {
	ctx        context.Context
	devices    devices
	batchSize  uint32
	symbol     []byte
	weights    []byte
	inputNodes []inputNode
	outputNode string
}

type Option func(*Options)

func Context(c context.Context) Option {
	return func(o *Options) {
		o.ctx = c
	}
}

func (o *Options) Context() context.Context {
	return o.ctx
}

func PredictorOptions(p dl.PredictionOptions) Option {
	return func(o *Options) {
		o.batchSize = p.BatchSize
	}
}

func (o *Options) PredictorOptions() dl.PredictionOptions {
	return dl.PredictionOptions{
		BatchSize: o.batchSize,
	}
}

func BatchSize(n uint32) Option {
	return func(o *Options) {
		o.batchSize = n
	}
}

func (o *Options) BatchSize() uint32 {
	return o.batchSize
}

func Device(deviceType DeviceType, id int) Option {
	return func(o *Options) {
		o.devices = append(o.devices, device{deviceType: deviceType, id: id})
	}
}

func (o *Options) Devices() devices {
	if len(o.devices) == 0 {
		return []device{Config.DefaultDevice}
	}
	return o.devices
}

func Symbol(sym []byte) Option {
	return func(o *Options) {
		o.symbol = sym
	}
}

func (o *Options) Symbol() []byte {
	return o.symbol
}

func Weights(w []byte) Option {
	return func(o *Options) {
		o.weights = w
	}
}

func (o *Options) Weights() []byte {
	return o.weights
}

func InputNode(key string, shape []uint32) Option {
	return func(o *Options) {
		o.inputNodes = append(
			o.inputNodes,
			inputNode{
				key:   key,
				shape: shape,
			},
		)
	}
}

func (o *Options) InputNodes() []inputNode {
	return o.inputNodes
}

func OutputNode(output string) Option {
	return func(o *Options) {
		o.outputNode = output
	}
}

func (o *Options) OutputNode() string {
	return o.outputNode
}

func New(opts ...Option) *Options {
	options := &Options{
		ctx:       context.Background(),
		batchSize: Config.BatchSize,
		devices:   []device{},
	}

	for _, o := range opts {
		o(options)
	}

	for ii, inputNode := range options.inputNodes {
		batchSize := uint32(options.batchSize)
		if len(options.inputNodes[ii].shape) == 3 {
			options.inputNodes[ii].shape = append([]uint32{batchSize}, inputNode.shape...)
		} else {
			options.inputNodes[ii].shape[0] = batchSize
		}
	}

	return options
}
