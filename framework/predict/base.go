package predict

import (
	"context"

	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/options"
	"github.com/rai-project/tracer"
)

type Base struct {
	Framework dlframework.FrameworkManifest
	Model     dlframework.ModelManifest
	Options   *options.Options
}

func (b Base) Info() (dlframework.FrameworkManifest, dlframework.ModelManifest, error) {
	return b.Framework, b.Model, nil
}

func (b Base) GetPredictionOptions(ctx context.Context) (*options.Options, error) {
	return b.Options, nil
}

func (b Base) BatchSize() int {
	return b.Options.BatchSize()
}

func (b Base) FeatureLimit() int {
	return b.Options.FeatureLimit()
}

func (b Base) TraceLevel() tracer.Level {
	return b.Options.TraceLevel()
}
