package common

import (
	"path/filepath"
	"sync"

	"github.com/Masterminds/semver"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/rai-project/dlframework"
	yaml "gopkg.in/yaml.v2"
)

func Register(framework dlframework.FrameworkManifest, a *assetfs.AssetFS) error {
	if err := framework.Register(); err != nil {
		log.WithField("framework", framework.GetName()).
			Error("failed to register framework")
		return err
	}
	frameworkVersion, err := semver.NewVersion(framework.GetVersion())
	if err != nil {
		log.WithField("frameworkVersion", framework.GetVersion()).
			Error("failed to parse framework version")
		return err
	}
	assets, err := a.AssetDir("")
	if err != nil {
		return err
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
				log.WithField("asset", asset).Error("failed to get asset bytes")
				return
			}

			var model dlframework.ModelManifest
			if err := yaml.Unmarshal(bts, &model); err != nil {
				log.WithField("asset", asset).Error("failed to unmarshal model ")
				return
			}
			if model.GetFramework().GetName() != framework.GetName() {
				log.WithField("asset", asset).Error("empty model name")
				return
			}
			modelFrameworkConstraint, err := semver.NewConstraint(model.GetFramework().GetVersion())
			if err != nil {
				log.WithField("modelFrameworkConstraint", model.GetFramework().GetVersion()).
					Error("failed to create model constraints")
				return
			}
			check := modelFrameworkConstraint.Check(frameworkVersion)
			if !check {
				log.WithField("frameworkVersion", frameworkVersion).
					WithField("modelFrameworkConstraint", model.GetFramework().GetVersion()).
					Error("failed to satisfy framework constraints")
				return
			}
			model.Register()
		}(ii, asset)
	}
	wg.Wait()
	return nil
}
