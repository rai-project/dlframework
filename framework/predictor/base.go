package predictor

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"

	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/options"
	"github.com/rai-project/tracer"
)

type Base struct {
	Framework dlframework.FrameworkManifest
	Model     dlframework.ModelManifest
	WorkDir   string
	Options   *options.Options
}

func (b Base) Info() (dlframework.FrameworkManifest, dlframework.ModelManifest, error) {
	return b.Framework, b.Model, nil
}

func (b Base) Modality() (dlframework.Modality, error) {
	return dlframework.UnknownModality, errors.New("undefined modality for predictor")
}

func (b Base) GetPredictionOptions(ctx context.Context) (*options.Options, error) {
	return b.Options, nil
}

func (b Base) BatchSize() int {
	return b.Options.BatchSize()
}

func (b Base) FeatureLimit() int {
	return b.Options.FeatureLimit()
}

func (b Base) TraceLevel() tracer.Level {
	return b.Options.TraceLevel()
}

func (b Base) UseGPU() bool {
	return b.Options.UsesGPU()
}

func (p Base) GetTypeParameter(typeParameters map[string]*dlframework.ModelManifest_Type_Parameter, name string) (string, error) {
	if typeParameters == nil {
		return "", errors.New("invalid type parameters")
	}
	pdims, ok := typeParameters[name]
	if !ok {
		return "", errors.New("expecting a type parameter")
	}
	pdimsVal := pdims.Value
	if pdimsVal == "" {
		return "", errors.New("invalid type parameter")
	}
	var ret string
	if err := yaml.Unmarshal([]byte(pdimsVal), &ret); err != nil {
		return "", errors.Errorf("unable to get the type parameter %v as a string", pdimsVal)
	}
	return ret, nil
}

func (p Base) GetPreprocessOptions(ctx context.Context) (PreprocessOptions, error) {
	return PreprocessOptions{}, errors.New("invalid preprocessor options")
}

func (p Base) baseURL(model dlframework.ModelManifest) string {
	baseURL := ""
	if model.GetModel().GetBaseUrl() != "" {
		baseURL = strings.TrimSuffix(model.GetModel().GetBaseUrl(), "/") + "/"
	}
	return baseURL
}

func (p Base) GetWeightsUrl() string {
	model := p.Model
	if model.GetModel().GetIsArchive() {
		return model.GetModel().GetBaseUrl()
	}
	if model.GetModel().GetWeightsPath() == "" {
		return ""
	}
	return strings.TrimRight(p.baseURL(model), "/") + "/" + model.GetModel().GetWeightsPath()
}

func (p Base) GetGraphUrl() string {
	model := p.Model
	if model.GetModel().GetIsArchive() {
		return model.GetModel().GetBaseUrl()
	}
	if model.GetModel().GetGraphPath() == "" {
		return ""
	}
	url := strings.TrimRight(p.baseURL(model), "/")
	if url != "" {
		url += "/"
	}
	return url + model.GetModel().GetGraphPath()
}

func (p Base) GetFeaturesUrl() string {
	model := p.Model
	params := model.GetOutput().GetParameters()
	pfeats, ok := params["features_url"]
	if !ok {
		return ""
	}
	return pfeats.Value
}

func (p Base) GetWeightsChecksum() string {
	model := p.Model
	return model.GetModel().GetWeightsChecksum()
}

func (p Base) GetGraphChecksum() string {
	model := p.Model
	return model.GetModel().GetGraphChecksum()
}

func (p Base) GetFeaturesChecksum() string {
	model := p.Model
	params := model.GetOutput().GetParameters()
	pfeats, ok := params["features_checksum"]
	if !ok {
		return ""
	}
	return pfeats.Value
}

func (p Base) GetWeightsPath() string {
	model := p.Model
	if model.GetModel().GetWeightsPath() == "" {
		return ""
	}
	graphPath := filepath.Base(model.GetModel().GetWeightsPath())
	return filepath.Join(p.WorkDir, graphPath)
}

func (p Base) GetGraphPath() string {
	model := p.Model
	if model.GetModel().GetGraphPath() == "" {
		return ""
	}
	graphPath := filepath.Base(model.GetModel().GetGraphPath())
	if graphPath == "" {
		return ""
	}
	return filepath.Join(p.WorkDir, graphPath)
}

func (p Base) GetFeaturesPath() string {
	model := p.Model
	return filepath.Join(p.WorkDir, model.GetName()+".features")
}

func (p Base) GetFeatureType() dlframework.FeatureType {
	model := p.Model
	ty := strings.ToUpper(model.GetOutput().GetType())
	return dlframework.FeatureType(dlframework.FeatureType_value[ty])
}
