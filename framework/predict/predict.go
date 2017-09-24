package predict

import (
	"io"

	"golang.org/x/net/context"

	"github.com/rai-project/dlframework"
)

type Predictor interface {
	// Gets framework and model manifests
	Info() (dlframework.FrameworkManifest, dlframework.ModelManifest, error)
	// Load model from manifest
	Load(ctx context.Context, model dlframework.ModelManifest, opts dlframework.PredictionOptions) (Predictor, error)
	// Returns the PredictionOptions options
	GetPredictionOptions(ctx context.Context) (dlframework.PredictionOptions, error)
	// Returns the preprocess options
	GetPreprocessOptions(ctx context.Context) (PreprocessOptions, error)
	// Returns the features
	Predict(ctx context.Context, data []float32, opts dlframework.PredictionOptions) (dlframework.Features, error)
	// Clears the internal state of a predictor
	Reset(ctx context.Context) error

	io.Closer
}
