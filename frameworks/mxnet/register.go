package mxnet

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/rai-project/dlframework"
	yaml "gopkg.in/yaml.v1"
)

var thisFramework = dlframework.FrameworkManifest{
	Name: "MXNet",
	Container: map[string]*dlframework.ContainerHardware{
		"amd64": &dlframework.ContainerHardware{
			Cpu: "raiproject/carml-mxnet:amd64-cpu",
			Gpu: "raiproject/carml-mxnet:amd64-gpu",
		},
		"ppc64le": &dlframework.ContainerHardware{
			Cpu: "raiproject/carml-mxnet:ppc64le-gpu",
			Gpu: "raiproject/carml-mxnet:ppc64le-gpu",
		},
	},
}

func init() {
	thisFramework.Register()

	assets, err := AssetDir("")
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

			bts, err := Asset(asset)
			if err != nil {
				return
			}

			var model dlframework.ModelManifest
			if err := yaml.Unmarshal(bts, &model); err != nil {
				return
			}
			name := strings.TrimRight(filepath.Base(asset), ext)
			model.Register()
		}(ii, asset)
	}
	wg.Wait()
}
