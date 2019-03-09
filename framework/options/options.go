package options

import (
	"context"
	"strings"

	dl "github.com/rai-project/dlframework"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/tracer"
)

type Options struct {
	ctx          context.Context
	devices      devices
	batchSize    int
	featureLimit int
	traceLevel   tracer.Level
	graph        []byte
	weights      []byte
}

type Option func(*Options)

func WithOptions(opts *Options) Option {
	return func(o *Options) {
		*o = *opts
	}
}

func Context(c context.Context) Option {
	return func(o *Options) {
		o.ctx = c
	}
}

func (o *Options) Context() context.Context {
	return o.ctx
}

func BatchSize(n int) Option {
	return func(o *Options) {
		o.batchSize = n
	}
}

func (o *Options) BatchSize() int {
	if o.batchSize == 0 {
		return 1
	}
	return o.batchSize
}

func FeatureLimit(num int) Option {
	return func(o *Options) {
		o.featureLimit = num
	}
}

func (o *Options) FeatureLimit() int {
	return o.featureLimit
}

func Device(deviceType DeviceType, id int) Option {
	return func(o *Options) {
		if deviceType == CUDA_DEVICE && !nvidiasmi.HasGPU {
			panic("cannot set CUDA device on systems with no GPU")
		}
		o.devices = append(o.devices, device{deviceType: deviceType, id: id})
	}
}

func (o *Options) Devices() devices {
	if len(o.devices) == 0 {
		return []device{Config.DefaultDevice}
	}
	return o.devices
}

func (o *Options) UsesGPU() bool {
	devs := o.Devices()
	for _, d := range devs {
		if d.Type() == CUDA_DEVICE {
			return true
		}
	}
	return false
}

func TraceLevel(tl dl.ExecutionOptions_TraceLevel) Option {
	return func(o *Options) {
		o.traceLevel = tracer.LevelFromName(tl.String())
	}
}

func (o *Options) TraceLevel() tracer.Level {
	return o.traceLevel
}

func Graph(g []byte) Option {
	return func(o *Options) {
		o.graph = g
	}
}

func (o *Options) Graph() []byte {
	return o.graph
}

func Weights(w []byte) Option {
	return func(o *Options) {
		o.weights = w
	}
}

func (o *Options) Weights() []byte {
	return o.weights
}

func (o *Options) Append(opts ...Option) *Options {
	for _, oi := range opts {
		oi(o)
	}
	return o
}

func PredictorOptions(p *dl.PredictionOptions) Option {
	return func(o *Options) {
		for k, v := range p.GetExecutionOptions().GetDeviceCount() {
			k = strings.ToLower(k)
			if k == "cpu" {
				o.devices = append(o.devices, device{deviceType: CPU_DEVICE, id: int(v)})
			} else {
				o.devices = append(o.devices, device{deviceType: CUDA_DEVICE, id: int(v)})
			}
		}
		o.batchSize = int(p.BatchSize)
		o.featureLimit = int(p.FeatureLimit)
		o.traceLevel = tracer.LevelFromName(p.GetExecutionOptions().GetTraceLevel().String())
	}
}

func New(opts ...Option) *Options {
	options := &Options{
		ctx:          context.Background(),
		devices:      []device{},
		batchSize:    Config.BatchSize,
		featureLimit: Config.FeatureLimit,
		traceLevel:   tracer.NO_TRACE,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}
