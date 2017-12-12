package predict

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	"golang.org/x/net/context"
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
	WorkDir string
}

func (p ImagePredictor) baseURL(model dlframework.ModelManifest) string {
	baseURL := ""
	if model.GetModel().GetBaseUrl() != "" {
		baseURL = strings.TrimSuffix(model.GetModel().GetBaseUrl(), "/") + "/"
	}
	return baseURL
}

func (p ImagePredictor) GetWeightsUrl() string {
	model := p.Model
	if model.GetModel().GetIsArchive() {
		return model.GetModel().GetBaseUrl()
	}
	return p.baseURL(model) + model.GetModel().GetWeightsPath()
}

func (p ImagePredictor) GetGraphUrl() string {
	model := p.Model
	if model.GetModel().GetIsArchive() {
		return model.GetModel().GetBaseUrl()
	}
	return p.baseURL(model) + model.GetModel().GetGraphPath()
}

func (p ImagePredictor) GetFeaturesUrl() string {
	model := p.Model
	params := model.GetOutput().GetParameters()
	pfeats, ok := params["features_url"]
	if !ok {
		return ""
	}
	return pfeats.Value
}

func (p ImagePredictor) GetWeightsChecksum() string {
	model := p.Model
	return model.GetModel().GetWeightsChecksum()
}

func (p ImagePredictor) GetGraphChecksum() string {
	model := p.Model
	return model.GetModel().GetGraphChecksum()
}

func (p ImagePredictor) GetFeaturesChecksum() string {
	model := p.Model
	params := model.GetOutput().GetParameters()
	pfeats, ok := params["features_checksum"]
	if !ok {
		return ""
	}
	return pfeats.Value
}

func (p ImagePredictor) GetWeightsPath() string {
	model := p.Model
	graphPath := filepath.Base(model.GetModel().GetWeightsPath())
	return cleanPath(filepath.Join(p.WorkDir, graphPath))
}

func (p ImagePredictor) GetGraphPath() string {
	model := p.Model
	graphPath := filepath.Base(model.GetModel().GetGraphPath())
	return cleanPath(filepath.Join(p.WorkDir, graphPath))
}

func (p ImagePredictor) GetFeaturesPath() string {
	model := p.Model
	return cleanPath(filepath.Join(p.WorkDir, model.GetName()+".features"))
}

func (p ImagePredictor) GetMeanPath() string {
	model := p.Model
	return cleanPath(filepath.Join(p.WorkDir, model.GetName()+".mean"))
}

func (p ImagePredictor) GetImageDimensions() ([]uint32, error) {
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

	var dims []uint32
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

func (p ImagePredictor) GetLayerName(typeParameters map[string]*dlframework.ModelManifest_Type_Parameter) (string, error) {
	if typeParameters == nil {
		return "", errors.New("invalid type parameters")
	}
	pdims, ok := typeParameters["layer_name"]
	if !ok {
		return "", errors.New("expecting a layer name")
	}
	pdimsVal := pdims.Value
	if pdimsVal == "" {
		return "", errors.New("invalid layer name")
	}

	var name string
	if err := yaml.Unmarshal([]byte(pdimsVal), &name); err != nil {
		return "", errors.Errorf("unable to get the layer name %v as a string", pdimsVal)
	}
	return name, nil
}

func (p ImagePredictor) GetInputLayerName(defaultValue string) string {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	name, err := p.GetLayerName(typeParameters)
	if err != nil {
		if defaultValue == "" {
			return DefaultInputLayerName
		}
		return defaultValue
	}
	return name
}

func (p ImagePredictor) GetOutputLayerName(defaultValue string) string {
	model := p.Model
	modelOutput := model.GetOutput()
	typeParameters := modelOutput.GetParameters()
	name, err := p.GetLayerName(typeParameters)
	if err != nil {
		if defaultValue == "" {
			return DefaultOutputLayerName
		}
		return defaultValue
	}
	return name
}

func (p ImagePredictor) GetPreprocessOptions(ctx context.Context) (PreprocessOptions, error) {
	return PreprocessOptions{}, errors.New("invalid preprocessor options")
}

func (p ImagePredictor) Reset(ctx context.Context) error {
	return nil
}

func (p ImagePredictor) Close() error {
	return nil
}
