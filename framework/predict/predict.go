package predict

import (
	"image"
	"io"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	yaml "gopkg.in/yaml.v2"

	"github.com/pkg/errors"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/utils"
)

type Predictor interface {
	// Gets framework information
	GetFramework() (dlframework.FrameworkManifest, error)
	// Load model from manifest
	Load(ctx context.Context, model dlframework.ModelManifest) (Predictor, error)
	// Downloads the features / symbol file / weights
	Download(ctx context.Context) error
	// Returns the features
	PredictURL(ctx context.Context, url string) (chan<- dlframework.PredictionFeature, error)
	// Returns the features
	PredictImage(ctx context.Context, image image.Image) (chan<- dlframework.PredictionFeatures, error)
	// Returns the features
	PredictDataset(ctx context.Context, dataset dldataset.Dataset) (chan<- dlframework.PredictionFeatures, error)
	// Preprocess image
	PreprocessImage(ctx context.Context, image image.Image) (image.Image, error)

	io.Closer
}

type Base struct {
	Framework dlframework.FrameworkManifest
	Model     dlframework.ModelManifest
}

type ImagePredictor struct {
	Base
	WorkDir string
}

func (p Base) GetFramework() (dlframework.FrameworkManifest, error) {
	return p.Framework, nil
}

func (p ImagePredictor) GetWeightsUrl() string {
	model := p.Model
	if model.GetModel().GetIsArchive() {
		return model.GetModel().GetBaseUrl()
	}
	baseURL := ""
	if model.GetModel().GetBaseUrl() != "" {
		baseURL = strings.TrimSuffix(model.GetModel().GetBaseUrl(), "/") + "/"
	}
	return baseURL + model.GetModel().GetWeightsPath()
}

func (p ImagePredictor) GetGraphUrl() string {
	model := p.Model
	if model.GetModel().GetIsArchive() {
		return model.GetModel().GetBaseUrl()
	}
	baseURL := ""
	if model.GetModel().GetBaseUrl() != "" {
		baseURL = strings.TrimSuffix(model.GetModel().GetBaseUrl(), "/") + "/"
	}
	return baseURL + model.GetModel().GetGraphPath()
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

func (p ImagePredictor) GetGraphPath() string {
	model := p.Model
	graphPath := filepath.Base(model.GetModel().GetGraphPath())
	return cleanPath(filepath.Join(p.WorkDir, graphPath))
}

func (p ImagePredictor) GetWeightsPath() string {
	model := p.Model
	graphPath := filepath.Base(model.GetModel().GetWeightsPath())
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
	pmean, ok := typeParameters["mean"]
	if !ok {
		log.Debug("using 0,0,0 as the mean image")
		return []float32{0, 0, 0}, nil
	}

	pmeanVal := pmean.Value
	if pmeanVal == "" {
		return nil, errors.New("invalid mean image")
	}

	if utils.IsURL(pmeanVal) {
		return urlProcessor(ctx, pmeanVal)
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
		return 1.0, errors.New("invalid type paramters")
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

func NoMeanImageURLProcessor(ctx context.Context, url string) ([]float32, error) {
	return nil, errors.New("mean image url processor disabled")
}

func cleanPath(path string) string {
	path = strings.Replace(path, ":", "_", -1)
	path = strings.Replace(path, " ", "_", -1)
	path = strings.Replace(path, "-", "_", -1)
	return strings.ToLower(path)
}
