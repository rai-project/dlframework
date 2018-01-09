package registryquery

import (
	"path"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/parallel/tunny"
	kv "github.com/rai-project/registry"
	"github.com/rai-project/serializer"
)

type frameworksTy struct {
	serializer serializer.Serializer
}

var Frameworks frameworksTy

func (f frameworksTy) Manifests() ([]*webmodels.DlframeworkFrameworkManifest, error) {
	rgs, err := kv.New()
	if err != nil {
		return nil, err
	}
	defer rgs.Close()

	var manifestsLock sync.Mutex
	var wg sync.WaitGroup
	manifests := []*webmodels.DlframeworkFrameworkManifest{}

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
			return errors.Errorf("invalid value for key=%s", e.Key)
		}

		framework := new(dlframework.FrameworkManifest)
		if err := f.serializer.Unmarshal(registryValue, framework); err != nil {
			return err
		}
		res := new(webmodels.DlframeworkFrameworkManifest)
		if err := copier.Copy(res, framework); err != nil {
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
	frameworksKey := path.Join(prefixKey, "frameworks")
	frameworksValue, err := rgs.Get(frameworksKey)
	if err != nil {
		return nil, err
	}
	frameworks, err := f.ProcessFrameworkNames(frameworksValue.Value)
	if err != nil {
		return nil, err
	}
	for _, framework := range frameworks {
		wg.Add(1)
		frameworkName, frameworkVersion := framework[0], framework[1]
		key := path.Join(prefixKey, frameworkName, frameworkVersion, "manifest.json")
		pool.SendWorkAsync(key, func(interface{}, error) {
			wg.Done()
		})
	}
	wg.Wait()
	return manifests, nil
}

func (frameworksTy) FilterManifests(
	manifests []*webmodels.DlframeworkFrameworkManifest,
	frameworkName,
	frameworkVersionString string,
) ([]*webmodels.DlframeworkFrameworkManifest, error) {
	frameworkName = strings.ToLower(frameworkName)
	frameworkVersionString = strings.ToLower(frameworkVersionString)

	candidates := []*webmodels.DlframeworkFrameworkManifest{}
	for _, manifest := range manifests {
		if frameworkName == "*" || strings.ToLower(manifest.Name) == frameworkName {
			candidates = append(candidates, manifest)
		}
	}
	if len(candidates) == 0 {
		return nil, errors.Errorf("framework %s not found", frameworkName)
	}

	if frameworkVersionString == "" || frameworkVersionString == "*" {
		return candidates, nil
	}

	sortByVersion := func(lst []*webmodels.DlframeworkFrameworkManifest) func(ii, jj int) bool {
		return func(ii, jj int) bool {
			f1, e1 := semver.NewVersion(lst[ii].Version)
			if e1 != nil {
				return false
			}
			f2, e2 := semver.NewVersion(lst[jj].Version)
			if e2 != nil {
				return false
			}
			return f1.LessThan(f2)
		}
	}

	if frameworkVersionString == "latest" {
		sort.Slice(candidates, sortByVersion(candidates))
		return []*webmodels.DlframeworkFrameworkManifest{candidates[0]}, nil
	}

	frameworkVersion, err := semver.NewConstraint(frameworkVersionString)
	if err != nil {
		return nil, err
	}

	res := []*webmodels.DlframeworkFrameworkManifest{}
	for _, manifest := range manifests {

		c, err := semver.NewVersion(manifest.Version)
		if err != nil {
			continue
		}
		if !frameworkVersion.Check(c) {
			continue
		}
		res = append(res, manifest)
	}
	if len(res) == 0 {
		return nil, errors.Errorf("framework %s=%s not found", frameworkName, frameworkVersionString)
	}
	sort.Slice(res, sortByVersion(res))

	return []*webmodels.DlframeworkFrameworkManifest{res[0]}, nil
}

func (f frameworksTy) ProcessFrameworkNames(buf []byte) ([][]string, error) {
	lines := strings.Split(string(buf), "\n")
	res := [][]string{}
	for _, line := range lines {
		res = append(res, strings.Split(line, ":"))
	}
	return res, nil
}

func init() {
	config.AfterInit(func() {
		Frameworks = frameworksTy{
			serializer: kv.Config.Serializer,
		}
	})
}
