package predict

import (
	"errors"
	"io"

	"github.com/gogo/protobuf/types"
	"github.com/rai-project/dlframework"
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
	if pdimsVal == nil {
		return nil, errors.New("invalid image dimensions")
	}
	data, ok := pdimsVal.Fields["data"]
	if !ok {
		return nil, errors.New("expecting data field in struct")
	}
	lstVal := data.GetListValue()
	if lstVal == nil {
		return nil, errors.New("expecting list value in data field in struct")
	}

	dims := []int32{}
	for _, v := range lstVal.Values {
		kind := v.GetKind()
		if kind == nil {
			return nil, errors.New("unable to get kind of value in image dimensions")
		}
		if _, ok := kind.(*types.Value_NumberValue); !ok {
			return nil, errors.New("invalid number value in image dimensions")
		}
		val := v.GetNumberValue()
		dims = append(dims, int32(val))
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
	if pdimsVal == nil {
		return nil, errors.New("invalid image dimensions")
	}
	data, ok := pdimsVal.Fields["data"]
	if !ok {
		return nil, errors.New("expecting data field in struct")
	}
	lstVal := data.GetListValue()
	if lstVal == nil {
		// try to get a number value
		kind := data.GetKind()
		if kind == nil {
			return nil, errors.New("unable to get kind of value in mean image")
		}
		if _, ok := kind.(*types.Value_NumberValue); !ok {
			return nil, errors.New("invalid number or list value in image mean")
		}
		val := float32(data.GetNumberValue())
		log.Debugf("using %v,%v,%v as the mean image", val, val, val)
		return []float32{val, val, val}, nil
	}

	dims := []float32{}
	for _, v := range lstVal.Values {
		kind := v.GetKind()
		if kind == nil {
			return nil, errors.New("unable to get kind of value in image dimensions")
		}
		if _, ok := kind.(*types.Value_NumberValue); !ok {
			return nil, errors.New("invalid number value in image dimensions")
		}
		val := v.GetNumberValue()
		dims = append(dims, float32(val))
	}
	return dims, nil
}
