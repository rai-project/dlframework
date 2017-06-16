package downloadmanager

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/hashicorp/go-getter"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

func cleanup(s string) string {
	return strings.Replace(s, ":", "_", -1)
}

func Download(url, targetPath string) error {

	// Get the string associated with the key url from the cache
	if _, found := cache.Get(url); found {
		return nil
	}

	targetPath = cleanup(targetPath)
	dirPath := filepath.Dir(targetPath)
	if !com.IsDir(dirPath) {
		err := os.MkdirAll(dirPath, 0700)
		if err != nil {
			return errors.Wrapf(err, "failed to create %v directory", dirPath)
		}
	}

	if com.IsFile(targetPath) {
		os.Remove(targetPath)
	}

	log.WithField("url", url).
		WithField("targetPath", targetPath).
		Debug("downloading data for prediction")

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
