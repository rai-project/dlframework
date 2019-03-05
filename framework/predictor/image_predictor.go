package predictor

import (
	"path/filepath"
	"sort"

	"context"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/feature"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	yaml "gopkg.in/yaml.v2"
)

type PreprocessOptions struct {
	Context   context.Context
	MeanImage []float32
	Size      []int
	Scale     float32
	ColorMode types.Mode
	Layout    image.Layout
}

type ImagePredictor struct {
	Base
	Metadata map[string]interface{}
}

func (p ImagePredictor) GetMeanPath() string {
	model := p.Model
	return cleanString(filepath.Join(p.WorkDir, model.GetName()+".mean"))
}

func (p ImagePredictor) GetImageDimensions() ([]int, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return nil, errors.New("invalid type parameters")
	}
	pdims, ok := typeParameters["dimensions"]
	if !ok {
		return nil, errors.New("expecting image type dimensions")
	}
	pdimsVal := pdims.Value
	if pdimsVal == "" {
		return nil, errors.New("invalid image dimensions")
	}

	var dims []int
	if err := yaml.Unmarshal([]byte(pdimsVal), &dims); err != nil {
		return nil, errors.Errorf("unable to get image dimensions %v as an integer slice", pdimsVal)
	}
	if len(dims) != 3 {
		return nil, errors.Errorf("expecting a dimensions size of 3, but got %v. do not put the batch size in the input dimensions.", len(dims))
	}
	return dims, nil
}

func (p ImagePredictor) GetMeanImage() ([]float32, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return nil, errors.New("invalid type parameters")
	}
	pmean, ok := typeParameters["mean"]
	if !ok {
		log.Debug("using 0,0,0 as the mean image")
		return []float32{0, 0, 0}, nil
	}

	pmeanVal := pmean.Value
	if pmeanVal == "" {
		return nil, errors.New("invalid mean image")
	}

	var vals []float32
	if err := yaml.Unmarshal([]byte(pmeanVal), &vals); err == nil {
		return vals, nil
	}
	var val float32
	if err := yaml.Unmarshal([]byte(pmeanVal), &val); err != nil {
		return nil, errors.Errorf("unable to get image mean %v as a float or slice", pmeanVal)
	}

	return []float32{val, val, val}, nil
}

func (p ImagePredictor) GetScale() (float32, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return 1.0, errors.New("invalid type parameters")
	}
	pscale, ok := typeParameters["scale"]
	if !ok {
		log.Debug("no scaling")
		return 1.0, nil
	}
	pscaleVal := pscale.Value
	if pscaleVal == "" {
		return 1.0, errors.New("invalid scale value")
	}

	var val float32
	if err := yaml.Unmarshal([]byte(pscaleVal), &val); err != nil {
		return 1.0, errors.Errorf("unable to get scale %v as a float", pscaleVal)
	}

	return val, nil
}

func (p ImagePredictor) GetLayout(defaultLayout image.Layout) image.Layout {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return defaultLayout
	}
	pscale, ok := typeParameters["layout"]
	if !ok {
		return defaultLayout
	}
	pscaleVal := pscale.Value
	if pscaleVal == "" {
		return defaultLayout
	}

	var val string
	if err := yaml.Unmarshal([]byte(pscaleVal), &val); err != nil {
		log.Errorf("unable to get color_mode %v as a string", pscaleVal)
		return defaultLayout
	}

	switch val {
	case "CHW":
		return image.CHWLayout
	case "HWC":
		return image.HWCLayout
	default:
		log.Error("invalid image mode specified " + val)
		return image.InvalidLayout
	}
}

func (p ImagePredictor) GetColorMode(defaultMode types.Mode) types.Mode {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return defaultMode
	}
	pscale, ok := typeParameters["color_mode"]
	if !ok {
		return defaultMode
	}
	pscaleVal := pscale.Value
	if pscaleVal == "" {
		return defaultMode
	}

	var val string
	if err := yaml.Unmarshal([]byte(pscaleVal), &val); err != nil {
		log.Errorf("unable to get color_mode %v as a string", pscaleVal)
		return defaultMode
	}

	switch val {
	case "RGB":
		return types.RGBMode
	case "BGR":
		return types.BGRMode
	default:
		log.Error("invalid image mode specified " + val)
		return types.InvalidMode
	}
}

// ReadPredictedFeatures ...
func (p ImagePredictor) CreateClassificationFeatures(ctx context.Context, probabilities []float32, labels []string) ([]dlframework.Features, error) {
	batchSize := p.BatchSize()
	featureLen := len(probabilities) / batchSize
	features := make([]dlframework.Features, batchSize)

	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		for jj := 0; jj < featureLen; jj++ {
			rprobs[jj] = feature.New(
				feature.ClassificationIndex(int32(jj)),
				feature.ClassificationLabel(labels[jj]),
				feature.Probability(probabilities[ii*featureLen+jj]),
			)
		}
		sort.Sort(dlframework.Features(rprobs))
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) CreateBoundingBoxFeatures(ctx context.Context, probabilities []float32, classes []float32, boxes [][]float32, labels []string) ([]dlframework.Features, error) {
	batchSize := p.BatchSize()
	featureLen := len(probabilities) / batchSize
	features := make([]dlframework.Features, batchSize)

	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		for jj := 0; jj < featureLen; jj++ {
			rprobs[jj] = feature.New(
				feature.BoundingBoxType(),
				feature.BoundingBoxXmin((boxes[ii*featureLen+jj][1])),
				feature.BoundingBoxXmax((boxes[ii*featureLen+jj][3])),
				feature.BoundingBoxYmin((boxes[ii*featureLen+jj][0])),
				feature.BoundingBoxYmax((boxes[ii*featureLen+jj][2])),
				feature.BoundingBoxLabel(labels[jj]),
				feature.BoundingBoxIndex(int32(classes[jj])),
				feature.Probability(probabilities[ii*featureLen+jj]),
			)
		}
		sort.Sort(dlframework.Features(rprobs))
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) Reset(ctx context.Context) error {
	return nil
}

func (p ImagePredictor) Close() error {
	return nil
}
