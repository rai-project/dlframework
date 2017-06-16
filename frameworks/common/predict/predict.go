package predict

import (
	"io"

	"github.com/rai-project/dlframework"
)

type Predictor interface {
	// Downloads the features / symbol file / weights
	Download() error
	// Preprocess the data
	Preprocess(data interface{}) (interface{}, error)
	// Returns the features
	Predict(data interface{}) ([]*dlframework.PredictionFeature, error)

	io.Closer
}

type Base struct {
	Framework dlframework.FrameworkManifest
	Model     dlframework.ModelManifest
}
