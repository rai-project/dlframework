package common

import (
	"path/filepath"
	"sync"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/rai-project/dlframework"
	yaml "gopkg.in/yaml.v2"
)

func Register(framework dlframework.FrameworkManifest, a *assetfs.AssetFS) {
	framework.Register()

	assets, err := a.AssetDir("")
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(assets))
	for ii, asset := range assets {
		go func(ii int, asset string) {
			defer wg.Done()
			ext := filepath.Ext(asset)
			if ext != ".yml" && ext != ".yaml" {
				return
			}

			bts, err := a.Asset(asset)
			if err != nil {
				return
			}

			var model dlframework.ModelManifest
			if err := yaml.Unmarshal(bts, &model); err != nil {
				return
			}
			model.Register()
		}(ii, asset)
	}
	wg.Wait()
}
