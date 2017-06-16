package downloadmanager

import (
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/hashicorp/go-getter"
	"github.com/k0kubun/pp"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
)

func cleanup(s string) string {
	return strings.Replace(s, ":", "_", -1)
}

func Download(url, targetDir string) (string, error) {

	// Get the string associated with the key url from the cache
	if val, found := cache.Get(url); found {
		s, ok := val.(string)
		if ok {
			return s, nil
		}
	}

	_, err := getter.Detect(url, targetDir, getter.Detectors)
	if err != nil {
		return "", err
	}

	targetDir = cleanup(targetDir)
	if !com.IsDir(targetDir) {
		err := os.MkdirAll(targetDir, 0700)
		if err != nil {
			return "", errors.Wrapf(err, "failed to create %v directory", targetDir)
		}
	}

	urlParsed, err := urlpkg.Parse(url)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse url %v", url)
	}
	filePath := filepath.Join(targetDir, filepath.Base(urlParsed.Path))
	if com.IsFile(filePath) || com.IsDir(filePath) {
		if config.IsDebug {
			log.Debugf("reusing the data in %v", filePath)
			return filePath, nil
		}
		os.RemoveAll(filePath)
	}
	pp.Println(filePath)

	log.WithField("url", url).
		WithField("targetDir", targetDir).
		Debug("downloading data for prediction")

	pwd := targetDir
	if com.IsFile(targetDir) {
		pwd = filepath.Dir(targetDir)
	}

	client := &getter.Client{
		Src:           url,
		Dst:           filePath,
		Pwd:           pwd,
		Mode:          getter.ClientModeFile,
		Decompressors: map[string]getter.Decompressor{}, // do not decompress
	}
	if err := client.Get(); err != nil {
		return "", err
	}

	unarchive(targetDir, filePath)

	// Set the value of the key url to targetDir, with the default expiration time
	cache.Set(url, filePath, gocache.DefaultExpiration)

	return filePath, nil
}

func unarchive(targetDir, filePath string) error {
	matchingLen := 0
	unArchiver := ""
	for k, _ := range getter.Decompressors {
		if strings.HasSuffix(filePath, "."+k) && len(k) > matchingLen {
			unArchiver = k
			matchingLen = len(k)
		}
	}
	if decompressor, ok := getter.Decompressors[unArchiver]; ok {
		decompressor.Decompress(targetDir, filePath, true)
	}
	return nil
}
