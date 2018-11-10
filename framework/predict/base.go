package predict

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/options"
	"github.com/rai-project/tracer"
	yaml "gopkg.in/yaml.v2"
)

type Base struct {
	Framework dlframework.FrameworkManifest
	Model     dlframework.ModelManifest
	Options   *options.Options
}

func (b Base) Info() (dlframework.FrameworkManifest, dlframework.ModelManifest, error) {
	return b.Framework, b.Model, nil
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

func (p Base) GetLayerName(typeParameters map[string]*dlframework.ModelManifest_Type_Parameter) (string, error) {
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

func (p Base) GetInputLayerName(defaultValue string) string {
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

func (p Base) GetOutputLayerName(defaultValue string) string {
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
	return p.baseURL(model) + model.GetModel().GetWeightsPath()
}

func (p Base) GetGraphUrl() string {
	model := p.Model
	if model.GetModel().GetIsArchive() {
		return model.GetModel().GetBaseUrl()
	}
	return p.baseURL(model) + model.GetModel().GetGraphPath()
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
	graphPath := filepath.Base(model.GetModel().GetWeightsPath())
	return cleanPath(filepath.Join(p.WorkDir, graphPath))
}

func (p Base) GetGraphPath() string {
	model := p.Model
	graphPath := filepath.Base(model.GetModel().GetGraphPath())
	return cleanPath(filepath.Join(p.WorkDir, graphPath))
}

func (p Base) GetFeaturesPath() string {
	model := p.Model
	return cleanPath(filepath.Join(p.WorkDir, model.GetName()+".features"))
}
