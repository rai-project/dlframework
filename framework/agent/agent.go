package agent

import (
	"path"
	"strings"
	"sync"

	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	dl "github.com/rai-project/dlframework"
	store "github.com/rai-project/libkv/store"
	lock "github.com/rai-project/lock/registry"
	"github.com/rai-project/registry"
)

type Base struct {
	Framework dl.FrameworkManifest
}

func toPath(s string) string {
	return strings.Replace(s, ":", "/", -1)
}

func (b *Base) PublishInRegistery(prefix string) error {
	var wg sync.WaitGroup

	framework := b.Framework
	cn, err := framework.CanonicalName()
	if err != nil {
		return err
	}

	ttl := registry.Config.Timeout
	marshaler := registry.Config.Serializer

	rgs, err := registry.New()
	if err != nil {
		return err
	}
	defer rgs.Close()

	prefix = path.Join(config.App.Name, prefix)

	rgs.Put(prefix, nil, &store.WriteOptions{IsDir: true})

	locker := lock.New(rgs)
	locker.Lock(prefix)
	defer locker.Unlock(prefix)

	frameworksKey := path.Join(prefix, "frameworks")

	kv, err := rgs.Get(frameworksKey)
	if err != nil {
		if ok, e := rgs.Exists(frameworksKey); e == nil && ok {
			log.WithError(err).Errorf("cannot get value for key %v", frameworksKey)
			return err
		}
		kv = &store.KVPair{
			Key:   frameworksKey,
			Value: []byte{},
		}
	}
	found := false
	val := strings.TrimSpace(string(kv.Value))
	frameworkLines := strings.Split(val, "\n")
	for _, name := range frameworkLines {
		if name == cn {
			found = true
			break
		}
	}
	if !found {
		frameworkLines = append(frameworkLines, cn)
		newVal := strings.TrimSpace(strings.Join(frameworkLines, "\n"))
		rgs.AtomicPut(frameworksKey, []byte(newVal), kv, nil)
	}

	key := path.Join(prefix, toPath(cn))
	rgs.Put(key, nil, &store.WriteOptions{TTL: ttl, IsDir: true})

	key = path.Join(key, "manifest.json")
	bts, err := marshaler.Marshal(&framework)
	if err != nil {
		return err
	}
	if err := rgs.Put(key, bts, &store.WriteOptions{TTL: ttl, IsDir: false}); err != nil {
		return err
	}

	models := framework.Models()
	wg.Add(len(models))
	for _, model := range models {
		go func(model dl.ModelManifest) {
			defer wg.Done()
			mn, err := model.CanonicalName()
			if err != nil {
				pp.Println(err)
				return
			}
			bts, err := marshaler.Marshal(&model)
			if err != nil {
				return
			}
			key := path.Join(prefix, toPath(mn), "manifest.json")
			rgs.Put(key, bts, &store.WriteOptions{TTL: ttl, IsDir: false})
		}(model)
	}

	wg.Wait()

	return nil
}
