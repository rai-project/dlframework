package registryquery

import (
	"encoding/json"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/jeffail/tunny"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	store "github.com/rai-project/libkv/store"
	kv "github.com/rai-project/registry"
)

func (m modelsTy) Agents(frameworkName, frameworkVersion, modelName, modelVersion string) ([]*webmodels.DlframeworkAgent, error) {

	frameworkName = strings.ToLower(frameworkName)
	frameworkVersion = strings.ToLower(frameworkVersion)
	modelName = strings.ToLower(modelName)
	modelVersion = strings.ToLower(modelVersion)

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

		pp.Println(val)

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
		frameworkName = strings.ToLower(model.Framework.Name)
		frameworkVersion = strings.ToLower(model.Framework.Version)
		modelName = strings.ToLower(model.Name)
		modelVersion = strings.ToLower(model.Version)

		// TODO:: the use of frameworkVersion here is not correct, since it won't support frameworkVersion=1.x.x for example
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
