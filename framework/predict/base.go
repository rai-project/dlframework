package predict

import (
	"github.com/rai-project/dlframework"
	"golang.org/x/net/context"
)

type Base struct {
	Framework         dlframework.FrameworkManifest
	Model             dlframework.ModelManifest
	PredictionOptions dlframework.PredictionOptions
}

func (b Base) GetPredictionOptions(ctx context.Context) (dlframework.PredictionOptions, error) {
	return b.PredictionOptions, nil
}

func (b Base) BatchSize() uint32 {
	s := b.PredictionOptions.GetBatchSize()
	if s == 0 {
		return uint32(1)
	}
	return s
}
