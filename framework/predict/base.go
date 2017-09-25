package predict

import (
	"github.com/rai-project/dlframework"
	tr "github.com/rai-project/tracer"
	"golang.org/x/net/context"
)

type Base struct {
	Framework         dlframework.FrameworkManifest
	Model             dlframework.ModelManifest
	PredictionOptions dlframework.PredictionOptions
	Tracer            tr.Tracer
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

func (b Base) GetTracer() tr.Tracer {
	return b.Tracer
}
