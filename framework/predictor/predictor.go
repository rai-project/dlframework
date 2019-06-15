package predictor

import (
	"io"
	"time"

	"context"

	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/options"
)

type Predictor interface {
	// Gets framework and model manifests
	Info() (dlframework.FrameworkManifest, dlframework.ModelManifest, error)
	// Gets predictor's Modality
	Modality() (dlframework.Modality, error)
	// Downloads model from manifest
	Download(ctx context.Context, model dlframework.ModelManifest, opts ...options.Option) (time.Duration, error)
	// Load model from manifest
	Load(ctx context.Context, model dlframework.ModelManifest, opts ...options.Option) (Predictor, time.Duration, error)
	// Returns the prediction options
	GetPredictionOptions() (*options.Options, error)
	// Returns the preprocess options
	GetPreprocessOptions() (PreprocessOptions, error)
	// Returns the handle to features
	Predict(ctx context.Context, data interface{}, opts ...options.Option) (time.Duration, time.Duration, error)
	// Returns the features
	ReadPredictedFeatures(ctx context.Context) ([]dlframework.Features, time.Duration, error)
	// Clears the internal state of a predictor
	Reset(ctx context.Context) error

	io.Closer
}
