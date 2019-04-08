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
	intputNodes  []Node
	outputNodes  []Node
}

type Option func(*Options)
type disableFrameworkAutoTuning struct{}

func WithOptions(opts *Options) Option {
	return func(o *Options) {
		*o = *opts
	}
}

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

func (o *Options) Context() context.Context {
	return o.ctx
}

func (o *Options) SetContext(ctx context.Context) {
	o.ctx = ctx
}

func DisableFrameworkAutoTuning(disabled bool) Option {
	return func(o *Options) {
		if o.ctx == nil {
			o.ctx = context.Background()
		}
		o.ctx = context.WithValue(o.ctx, disableFrameworkAutoTuning{}, disabled)
	}
}

func (o *Options) DisableFrameworkAutoTuning() bool {
	ctx := o.ctx
	if ctx == nil {
		return false
	}
	val, ok := ctx.Value(disableFrameworkAutoTuning{}).(bool)
	if !ok {
		return false
	}
	return val
}

func (o *Options) SetDisableFrameworkAutoTuning(disabled bool) {
	o.ctx = context.WithValue(o.ctx, disableFrameworkAutoTuning{}, disabled)
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

func (o *Options) SetBatchSize(n int) {
	o.batchSize = n
}

func FeatureLimit(num int) Option {
	return func(o *Options) {
		o.featureLimit = num
	}
}

func (o *Options) FeatureLimit() int {
	return o.featureLimit
}

func (o *Options) SetFeatureLimit(n int) {
	o.featureLimit = n
}

func Device(deviceType DeviceType, id int) Option {
	return func(o *Options) {
		if deviceType == CUDA_DEVICE && !nvidiasmi.HasGPU {
			panic("cannot set CUDA device on systems with no GPU")
		}
		o.devices = append(o.devices, device{deviceType: deviceType, id: id})
	}
}

func (o *Options) SetDevice(deviceType DeviceType, id int) {
	o.devices = []device{
		device{deviceType: deviceType, id: id},
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

func (o *Options) SetTraceLevel(g tracer.Level) {
	o.traceLevel = g
}

func Graph(g []byte) Option {
	return func(o *Options) {
		o.graph = g
	}
}

func (o *Options) Graph() []byte {
	return o.graph
}

func (o *Options) SetGraph(g []byte) {
	o.graph = g
}

func Weights(w []byte) Option {
	return func(o *Options) {
		o.weights = w
	}
}

func (o *Options) Weights() []byte {
	return o.weights
}

func (o *Options) SetWeights(w []byte) {
	o.weights = w
}

func (o *Options) Append(opts ...Option) *Options {
	for _, oi := range opts {
		oi(o)
	}
	return o
}

func InputNodes(ins []Node) Option {
	return func(o *Options) {
		o.intputNodes = ins
	}
}

func (o *Options) InputNodes() []Node {
	return o.intputNodes
}

func (o *Options) SetInputNodes(ins []Node) {
	o.intputNodes = ins
}

func OutputNodes(outs []Node) Option {
	return func(o *Options) {
		o.outputNodes = outs
	}
}

func (o *Options) OutputNodes() []Node {
	return o.outputNodes
}

func (o *Options) SetOutputNodes(ins []Node) {
	o.outputNodes = ins
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
