package predict

import (
	"io"

	"golang.org/x/net/context"
	yaml "gopkg.in/yaml.v2"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/utils"
)

type Predictor interface {
	// Downloads the features / symbol file / weights
	Download(ctx context.Context) error
	// Preprocess the data
	Preprocess(ctx context.Context, data interface{}) (interface{}, error)
	// Returns the features
	Predict(ctx context.Context, data interface{}) (*dlframework.PredictionFeatures, error)

	io.Closer
}

type Base struct {
	Framework dlframework.FrameworkManifest
	Model     dlframework.ModelManifest
}

type ImagePredictor struct {
	Base
}

func (p ImagePredictor) GetImageDimensions() ([]int32, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return nil, errors.New("invalid type paramters")
	}
	pdims, ok := typeParameters["dimensions"]
	if !ok {
		return nil, errors.New("expecting image type dimensions")
	}
	pdimsVal := pdims.Value
	if pdimsVal == "" {
		return nil, errors.New("invalid image dimensions")
	}

	var dims []int32
	if err := yaml.Unmarshal([]byte(pdimsVal), &dims); err != nil {
		return nil, errors.Errorf("unable to get image dimensions %v as an integer slice", pdimsVal)
	}
	return dims, nil
}

func NoMeanImageURLProcessor(ctx context.Context, url string) ([]float32, error) {
	return nil, errors.New("mean image url processor disabled")
}

func (p ImagePredictor) GetMeanImage(
	ctx context.Context,
	urlProcessor func(ctx context.Context, url string) ([]float32, error),
) ([]float32, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return nil, errors.New("invalid type paramters")
	}
	pdims, ok := typeParameters["mean"]
	if !ok {
		log.Debug("using 0,0,0 as the mean image")
		return []float32{0, 0, 0}, nil
	}
	pdimsVal := pdims.Value
	if pdimsVal == "" {
		return nil, errors.New("invalid image dimensions")
	}
	if utils.IsURL(pdimsVal) {
		return urlProcessor(ctx, pdimsVal)
	}

	var vals []float32
	if err := yaml.Unmarshal([]byte(pdimsVal), &vals); err == nil {
		return vals, nil
	}
	var val float32
	if err := yaml.Unmarshal([]byte(pdimsVal), &val); err != nil {
		return nil, errors.Errorf("unable to get image mean %v as a float or slice", pdimsVal)
	}

	return []float32{val, val, val}, nil
}
