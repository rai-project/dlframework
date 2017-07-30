package http

import (
	"path"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/parallel"
	kv "github.com/rai-project/registry"
	"github.com/rai-project/serializer"
)

type frameworksTy struct {
	serializer serializer.Serializer
}

var frameworks = frameworksTy{
	serializer: kv.Config.Serializer,
}

func (f frameworksTy) manifests() ([]*models.DlframeworkFrameworkManifest, error) {
	rgs, err := kv.New()
	if err != nil {
		return nil, err
	}
	defer rgs.Close()

	var manifestsLock sync.Mutex
	manifests := []*models.DlframeworkFrameworkManifest{}

	pr := parallel.New(0)

	dirs := []string{path.Join(config.App.Name, "registry")}
	for {
		if len(dirs) == 0 {
			break
		}
		var dir string
		dir, dirs = dirs[0], dirs[1:]
		lst, err := rgs.List(dir)
		if err != nil {
			continue
		}

		for ii := range lst {
			e := lst[ii]
			if e.Value == nil || len(e.Value) == 0 {
				dirs = append(dirs, e.Key)
				continue
			}
			pr.Add(parallel.NonCancelableTaskFunc(func() {
				registryValue := e.Value
				framework := new(dlframework.FrameworkManifest)
				if err := f.serializer.Unmarshal(registryValue, framework); err != nil {
					return
				}
				res := new(models.DlframeworkFrameworkManifest)
				if err := copier.Copy(res, framework); err != nil {
					return
				}
				manifestsLock.Lock()
				defer manifestsLock.Unlock()
				manifests = append(manifests, res)
			}))
		}
	}
	pr.Run()
	return manifests, nil
}
