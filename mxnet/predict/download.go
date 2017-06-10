package predict

import (
	"time"

	"github.com/levigross/grequests"
	gocache "github.com/patrickmn/go-cache"
)

var (
	cache *gocache.Cache
)

func downloadURL(url, targetPath string) error {

	// Get the string associated with the key url from the cache
	if _, found := cache.Get(url); found {
		return nil
	}

	log.WithField("url", url).WithField("targetPath", targetPath).Debug("downloading data for prediction")

	resp, err := grequests.Get(url, nil)
	if err != nil {
		return err
	}
	defer resp.Close()

	err = resp.DownloadToFile(targetPath)
	if err != nil {
		return err
	}

	// Set the value of the key url to targetPath, with the default expiration time
	cache.Set(url, targetPath, gocache.DefaultExpiration)

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

func init() {

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	cache = gocache.New(5*time.Minute, 10*time.Minute)

}
