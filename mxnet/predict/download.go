package predict

import (
	"io"
	"log"
	"net/http"
	"os"
)

func downloadURL(url, targetPath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	return nil
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
