package downloadmanager

import (
	"github.com/Unknwon/com"
	"github.com/hashicorp/go-getter"
	gocache "github.com/patrickmn/go-cache"
)

func Download(url, targetPath string) error {

	// Get the string associated with the key url from the cache
	if _, found := cache.Get(url); found {
		return nil
	}

	if com.IsFile(targetPath) {
		return nil
	}

	log.WithField("url", url).WithField("targetPath", targetPath).Debug("downloading data for prediction")

	client := &getter.Client{
		Src:  url,
		Dst:  targetPath,
		Pwd:  targetPath,
		Mode: getter.ClientModeAny,
	}

	err := client.Get()
	if err != nil {
		return err
	}

	// Set the value of the key url to targetPath, with the default expiration time
	cache.Set(url, targetPath, gocache.DefaultExpiration)

	return nil
}
