package predict

import (
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/options"
	"golang.org/x/net/context"
)

type Base struct {
	Framework dlframework.FrameworkManifest
	Model     dlframework.ModelManifest
	Options   *options.Options
}

func (b Base) GetPredictionOptions(ctx context.Context) (*options.Options, error) {
	return b.Options, nil
}

func (b Base) BatchSize() uint32 {
	return b.Options.BatchSize()
}
