package options

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type optionsConfig struct {
	BatchSize uint          `json:"batchsize" config:"predictor.batch_size"`
	done      chan struct{} `json:"-" config:"-"`
}

var (
	Config = &optionsConfig{
		done: make(chan struct{}),
	}
)

func (optionsConfig) ConfigName() string {
	return "predictor/options"
}

func (a *optionsConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *optionsConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	if a.BatchSize == 0 {
		a.BatchSize = DefaultBatchSize
	}
}

func (c optionsConfig) Wait() {
	<-c.done
}

func (c optionsConfig) String() string {
	return pp.Sprintln(c)
}

func (c optionsConfig) Debug() {
	log.Debug("predictor/options Config = ", c)
}

func init() {
	config.Register(Config)
}
