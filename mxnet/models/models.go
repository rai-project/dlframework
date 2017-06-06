package models

import (
	"path/filepath"
	"strings"
	"sync"

	yaml "gopkg.in/yaml.v1"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/mxnet"
)

type r struct {
	sync.RWMutex
	data map[string]mxnet.ModelInformation
}

var registry = r{
	data: map[string]mxnet.ModelInformation{},
}

func Names() []string {
	registry.RLock()
	defer registry.RUnlock()

	ii := 0
	names := make([]string, len(registry.data))
	for name := range registry.data {
		names[ii] = name
		ii++
	}

	return names
}

func Get(name string) (mxnet.ModelInformation, error) {
	registry.RLock()
	defer registry.RUnlock()

	name = strings.ToLower(name)
	if model, ok := registry.data[name]; ok {
		return model, nil
	}

	return mxnet.ModelInformation{}, errors.Errorf("cannot find model %s in registry", name)
}

func Register(model mxnet.ModelInformation) {
	registry.Lock()
	defer registry.Unlock()

	name := strings.ToLower(model.Name)
	registry.data[name] = model
}

func init() {
	assets, err := AssetDir("")
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(assets))
	for _, asset := range assets {
		go func(asset string) {
			defer wg.Done()
			ext := filepath.Ext(asset)
			if ext != ".yml" && ext != ".yaml" {
				return
			}

			bts, err := Asset(asset)
			if err != nil {
				return
			}

			var model mxnet.ModelInformation
			if err := yaml.Unmarshal(bts, &model); err != nil {
				return
			}

			Register(model)
		}(asset)
	}
	wg.Wait()
}
