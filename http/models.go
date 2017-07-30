package http

import (
	"path"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/jeffail/tunny"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	kv "github.com/rai-project/registry"
	"github.com/rai-project/serializer"
)

type modelsTy struct {
	serializer serializer.Serializer
}

var models modelsTy

func (m modelsTy) manifests(frameworkName, frameworkVersion string) ([]*webmodels.DlframeworkModelManifest, error) {

	frameworkName = strings.ToLower(frameworkName)
	frameworkVersion = strings.ToLower(frameworkVersion)

	fs, err := frameworks.manifests()
	if err != nil {
		return nil, err
	}
	fs, err = frameworks.filter(fs, frameworkName, frameworkVersion)
	rgs, err := kv.New()
	if err != nil {
		return nil, err
	}
	defer rgs.Close()

	var manifestsLock sync.Mutex
	var wg sync.WaitGroup
	manifests := []*webmodels.DlframeworkModelManifest{}

	poolSize := runtime.NumCPU()
	pool, err := tunny.CreatePool(poolSize, func(object interface{}) interface{} {
		key, ok := object.(string)
		if !ok {
			return errors.New("invalid key type. expecting a string type")
		}

		e, err := rgs.Get(key)
		if err != nil {
			return err
		}
		registryValue := e.Value
		if registryValue == nil || len(registryValue) == 0 {
			return nil
		}

		model := new(dlframework.ModelManifest)
		if err := m.serializer.Unmarshal(registryValue, model); err != nil {
			return err
		}
		res := new(webmodels.DlframeworkModelManifest)
		if err := copier.Copy(res, model); err != nil {
			return err
		}

		manifestsLock.Lock()
		defer manifestsLock.Unlock()

		manifests = append(manifests, res)
		return nil
	}).Open()
	if err != nil {
		return nil, err
	}

	defer pool.Close()

	prefixKey := path.Join(config.App.Name, "registry")
	for _, framework := range fs {
		frameworkName, frameworkVersion := strings.ToLower(framework.Name), strings.ToLower(framework.Version)
		key := path.Join(prefixKey, frameworkName, frameworkVersion)
		kvs, err := rgs.List(key)
		if err != nil {
			continue
		}
		for _, kv := range kvs {
			if path.Dir(kv.Key) == key {
				continue
			}
			wg.Add(1)
			pool.SendWorkAsync(kv.Key, func(interface{}, error) {
				wg.Done()
			})
		}
	}

	wg.Wait()
	return manifests, nil
}

func (m modelsTy) allmanifests() ([]*webmodels.DlframeworkModelManifest, error) {
	return m.manifests("*", "*")
}

func (modelsTy) filter(
	manifests []*webmodels.DlframeworkModelManifest,
	modelName,
	modelVersionString string,
) ([]*webmodels.DlframeworkModelManifest, error) {
	modelName = strings.ToLower(modelName)
	modelVersionString = strings.ToLower(modelVersionString)

	candidates := []*webmodels.DlframeworkModelManifest{}
	for _, manifest := range manifests {
		if modelName == "*" || strings.ToLower(manifest.Name) == modelName {
			candidates = append(candidates, manifest)
		}
	}
	if len(candidates) == 0 {
		return nil, errors.Errorf("model %s not found", modelName)
	}

	if modelVersionString == "" || modelVersionString == "*" {
		return candidates, nil
	}

	sortByVersion := func(ii, jj int) bool {
		f1, e1 := semver.NewVersion(candidates[ii].Version)
		if e1 != nil {
			return false
		}
		f2, e2 := semver.NewVersion(candidates[jj].Version)
		if e2 != nil {
			return false
		}
		return f1.LessThan(f2)
	}

	if modelVersionString == "latest" {
		sort.Slice(candidates, sortByVersion)
		return []*webmodels.DlframeworkModelManifest{candidates[0]}, nil
	}

	modelVersion, err := semver.NewConstraint(modelVersionString)
	if err != nil {
		return nil, err
	}

	res := []*webmodels.DlframeworkModelManifest{}
	for _, manifest := range manifests {

		c, err := semver.NewVersion(manifest.Version)
		if err != nil {
			continue
		}
		if !modelVersion.Check(c) {
			continue
		}
		res = append(res, manifest)
	}
	if len(res) == 0 {
		return nil, errors.Errorf("model %s=%s not found", modelName, modelVersionString)
	}
	sort.Slice(res, sortByVersion)

	return []*webmodels.DlframeworkModelManifest{res[0]}, nil
}

func init() {
	config.AfterInit(func() {
		models = modelsTy{
			serializer: kv.Config.Serializer,
		}
	})
}
