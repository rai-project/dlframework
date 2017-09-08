package predict

import (
	"io"

	"golang.org/x/net/context"

	"github.com/rai-project/dlframework"
	"github.com/rai-project/image/types"
)

type PreprocessOptions struct {
	MeanImage  []float32
	Size       []int
	Scale      int
	ColorSpace types.Mode
}

type Predictor interface {
	// Gets framework and model manifests
	Info() (dlframework.FrameworkManifest, dlframework.ModelManifest, error)
	// Load model from manifest
	Load(ctx context.Context, model dlframework.ModelManifest) (Predictor, error)
	// Returns the preprocessor options
	PreprocessOptions(ctx context.Context) (PreprocessOptions, error)
	// Returns the features
	Predict(ctx context.Context, data []float32) (dlframework.Features, error)
	// Clears the internal state of a predictor
	Reset(ctx context.Context) error

	io.Closer
}
