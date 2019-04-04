package registryquery

import (
	"encoding/json"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/Masterminds/semver"

	"github.com/pkg/errors"
	"github.com/rai-project/config"
	dl "github.com/rai-project/dlframework"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	store "github.com/rai-project/libkv/store"
	"github.com/rai-project/parallel/tunny"
	kv "github.com/rai-project/registry"
)

func (m modelsTy) Agents(frameworkName, frameworkVersion, modelName, modelVersion string) ([]*webmodels.DlframeworkAgent, error) {
	frameworkName = dl.CleanString(frameworkName)
	frameworkVersion = dl.CleanString(frameworkVersion)
	modelName = dl.CleanString(modelName)
	modelVersion = dl.CleanString(modelVersion)

	manifests, err := Models.Manifests(frameworkName, frameworkVersion)
	if err != nil {
		return nil, err
	}

	if len(manifests) == 0 {
		return nil, errors.Errorf("no models found for the framework %s:%s", frameworkName, frameworkVersion)
	}

	manifests, err = Models.FilterManifests(manifests, modelName, modelVersion)
	if err != nil {
		return nil, err
	}

	rgs, err := kv.New()
	if err != nil {
		return nil, err
	}
	defer rgs.Close()

	var agentsLock sync.Mutex
	var wg sync.WaitGroup
	set := make(map[string]bool)
	agents := []*webmodels.DlframeworkAgent{}

	poolSize := runtime.NumCPU()
	pool, err := tunny.CreatePool(poolSize, func(object interface{}) interface{} {
		kvs, ok := object.(*store.KVPair)
		if !ok {
			return errors.New("invalid kv type. expecting a KVPair type")
		}
		key := kvs.Key
		val := kvs.Value

		keyBase := path.Base(key)
		if !strings.HasPrefix(keyBase, "agent-") {
			return errors.Errorf("skipping non agent %s", keyBase)
		}

		hostPort := strings.Split(strings.TrimPrefix(keyBase, "agent-"), ":")
		host, port := hostPort[0], hostPort[1]

		agentsLock.Lock()
		defer agentsLock.Unlock()

		if _, ok := set[keyBase]; ok {
			return nil
		}

		agent := &webmodels.DlframeworkAgent{}
		err := json.Unmarshal(val, agent)
		if err != nil {
			log.WithError(err).WithField("host", host).WithField("port", port).Error("failed to unmarshal agent")
			return nil
		}

		agents = append(agents, agent)

		set[keyBase] = true
		return nil
	}).Open()
	if err != nil {
		return nil, err
	}

	defer pool.Close()

	prefixKey := path.Join(config.App.Name, "predictor")
	for _, model := range manifests {
		frameworkName = dl.CleanString(model.Framework.Name)
		frameworkVersion = dl.CleanString(model.Framework.Version)
		modelName = dl.CleanString(model.Name)
		modelVersion = dl.CleanString(model.Version)

		frameworkSemanticVersion, err := semver.NewVersion(frameworkVersion)
		if err != nil {
			log.WithError(err).Errorf("unable to get semantic version of %v", frameworkVersion)
			continue
		}

		frameworkKey := path.Join(prefixKey, frameworkName)
		frameworkEntries, err := rgs.List(frameworkKey)
		if err != nil {
			continue
		}
		for _, key := range frameworkEntries {
			fullPath := strings.Trim(key.Key, frameworkKey+"/")
			idx := strings.Index(fullPath, "/")
			registeredFrameworkVersion := strings.TrimSuffix(fullPath[:idx+1], "/")

			registeredFrameworkSemanticVersion, err := semver.NewVersion(registeredFrameworkVersion)
			if err != nil {
				continue
			}
			if frameworkSemanticVersion.Equal(registeredFrameworkSemanticVersion) {
				frameworkVersion = registeredFrameworkVersion
				break
			}
		}

		// TODO:: the use of modelVersion here is not correct, since it won't support modelVersion=1.x.x for example
		// same logic as being performed for framework version would suffice here
		key := path.Join(prefixKey, frameworkName, frameworkVersion, modelName, modelVersion)

		kvs, err := rgs.List(key)
		if err != nil {
			continue
		}
		for _, kv := range kvs {
			wg.Add(1)
			pool.SendWorkAsync(kv, func(interface{}, error) {
				wg.Done()
			})
		}
	}

	wg.Wait()

	return agents, nil
}
