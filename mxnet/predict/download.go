package predict

import (
	"github.com/pkg/errors"
)

func downloadURL(url, targetPath string) error {
	return errors.Errorf("implement me %s", url)
}

func (p *ImagePredictor) Download() error {
	model := p.model
	if model.GetGraphUrl() != "" {
		err := downloadURL(model.GetGraphUrl(), p.GetGraphPath())
		if err != nil {
			return err
		}
	}
	if model.GetWeightsUrl() != "" {
		err := downloadURL(model.GetWeightsUrl(), p.GetWeightsPath())
		if err != nil {
			return err
		}
	}
	if model.GetFeaturesUrl() != "" {
		err := downloadURL(model.GetFeaturesUrl(), p.GetFeaturesPath())
		if err != nil {
			return err
		}
	}
	return nil
}
