package framework

import (
	"path/filepath"

	"github.com/Masterminds/semver"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	yaml "gopkg.in/yaml.v2"
)

func Register(framework dlframework.FrameworkManifest, a *assetfs.AssetFS) error {
	if err := framework.Register(); err != nil {
		log.WithField("framework", framework.MustCanonicalName()).
			Error("failed to register framework")
		return err
	}
	if debugging {
		log.WithField("framework", framework.GetName()).
			Debug("registered framework")
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
	for _, asset := range assets {
		ext := filepath.Ext(asset)
		if ext != ".yml" && ext != ".yaml" {
			return err
		}

		bts, err := a.Asset(asset)
		if err != nil {
			log.WithField("asset", asset).Error("failed to get asset bytes")
			return err
		}

		var model dlframework.ModelManifest
		if err := yaml.Unmarshal(bts, &model); err != nil {
			log.WithField("asset", asset).WithError(err).Error("failed to unmarshal model")
			return err
		}
		if model.GetName() == "" {
			log.WithField("asset", asset).WithField("name", model.GetName()).Error("empty model name")
			return errors.New("empty model name found")
		}
		if model.GetFramework().GetName() != framework.GetName() {
			log.WithField("asset", asset).
				WithField("model_framework_name", model.GetFramework().GetName()).
				WithField("framework_name", framework.GetName()).
				Error("empty framework name")
			return err
		}
		modelFrameworkConstraint, err := semver.NewConstraint(model.GetFramework().GetVersion())
		if err != nil {
			log.WithField("model_framework_constraint", model.GetFramework().GetVersion()).
				Error("failed to create model constraints")
			return err
		}
		check := modelFrameworkConstraint.Check(frameworkVersion)
		if !check {
			log.WithField("frameworkVersion", frameworkVersion).
				WithField("model_framework_constraint", model.GetFramework().GetVersion()).
				Error("failed to satisfy framework constraints")
			return err
		}

		if model.GetHidden() {
			log.WithField("name", model.GetName()).Info("skipping regitration of hidden model")
			continue
		}

		if err := model.Register(); err != nil {
			log.WithError(err).
				WithField("frameworkVersion", frameworkVersion).
				WithField("model_name", model.GetName()).
				WithField("model_framework_constraint", model.GetFramework().GetVersion()).
				Error("failed to register model")
			continue
		}
		if debugging {
			log.WithField("framework", framework.MustCanonicalName()).
				WithField("model", model.MustCanonicalName()).
				Debug("registered model")
		}
	}
	return nil
}
