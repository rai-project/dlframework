package predict

import (
	"github.com/levigross/grequests"
)

func downloadURL(url, targetPath string) error {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		return err
	}
	defer resp.Close()
	return resp.DownloadToFile(targetPath)
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
