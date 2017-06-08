package mxnet

import (
	"path/filepath"
	"strings"
	"sync"

	yaml "gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

type modelInfoRegistryTy struct {
	sync.RWMutex
	data map[string]Model_Information
}

var modelInfoRegistry = modelInfoRegistryTy{
	data: map[string]Model_Information{},
}

func ModelNames() []string {
	modelInfoRegistry.RLock()
	defer modelInfoRegistry.RUnlock()

	ii := 0
	names := make([]string, len(modelInfoRegistry.data))
	for name := range modelInfoRegistry.data {
		names[ii] = name
		ii++
	}

	return names
}

func GetModelInformation(name string) (Model_Information, error) {
	modelInfoRegistry.RLock()
	defer modelInfoRegistry.RUnlock()

	name = strings.ToLower(name)
	if model, ok := modelInfoRegistry.data[name]; ok {
		return model, nil
	}

	return Model_Information{}, errors.Errorf("cannot find model %s in modelInfoRegistry", name)
}

func RegisterModelInformation(name string, model Model_Information) {
	modelInfoRegistry.Lock()
	defer modelInfoRegistry.Unlock()
	if name == "" {
		name = model.Name
	}
	name = strings.ToLower(name)
	modelInfoRegistry.data[name] = model
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

			var model Model_Information
			if err := yaml.Unmarshal(bts, &model); err != nil {
				return
			}

			name := strings.TrimRight(filepath.Base(asset), ext)
			RegisterModelInformation(name, model)

		}(asset)
	}
	wg.Wait()
}
