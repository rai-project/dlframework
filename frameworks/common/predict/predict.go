package predict

import (
	"io"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/spf13/cast"
)

type Predictor interface {
	// Downloads the features / symbol file / weights
	Download() error
	// Preprocess the data
	Preprocess(data interface{}) (interface{}, error)
	// Returns the features
	Predict(data interface{}) (*dlframework.PredictionFeatures, error)

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

	slice, err := cast.ToSliceE(pdimsVal)
	if err != nil {
		return nil, errors.Errorf("unable to get image dimensions %v as an integer slice", pdimsVal)
	}

	dims := []int32{}
	for _, v := range slice {
		val, err := cast.ToInt32E(v)
		if err != nil {
			return nil, errors.Errorf("unable to get image mean %v as an integer slice", pdimsVal)
		}
		dims = append(dims, val)
	}
	return dims, nil
}

func (p ImagePredictor) GetMeanImage() ([]float32, error) {
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

	slice, err := cast.ToSliceE(pdimsVal)
	if err != nil {
		val, err := cast.ToFloat32E(pdimsVal)
		if err != nil {
			return nil, errors.Errorf("unable to get image mean %v as a float or slice", pdimsVal)
		}
		log.Debugf("using %v,%v,%v as the mean image", val, val, val)
		return []float32{val, val, val}, nil
	}

	dims := []float32{}
	for _, v := range slice {
		f, err := cast.ToFloat32E(v)
		if err != nil {
			return nil, errors.Errorf("unable to get image mean %v as a float or slice", pdimsVal)
		}
		dims = append(dims, f)
	}
	return dims, nil
}
