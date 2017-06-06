package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	yaml "gopkg.in/yaml.v1"

	rice "github.com/GeertJohan/go.rice"
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
	var wg sync.WaitGroup
	var builtinBox = rice.MustFindBox("builtin")
	builtinBox.Walk(".", filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		wg.Add(1)
		go func() error {
			defer wg.Done()
			if err != nil {
				return fmt.Errorf("error walking box: %s\n", err)
			}

			ext := filepath.Ext(path)
			if ext != ".yml" && ext != ".yaml" {
				return nil
			}

			bts, err := builtinBox.Bytes(path)
			if err != nil {
				return err
			}

			var model mxnet.ModelInformation
			if err := yaml.Unmarshal(bts, &model); err != nil {
				return err
			}

			Register(model)

			return nil
		}()
		wg.Wait()
		return nil
	}))
}
