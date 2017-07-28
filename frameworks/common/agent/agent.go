package agent

import (
	"bytes"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	dl "github.com/rai-project/dlframework"
	store "github.com/rai-project/libkv/store"
	"github.com/rai-project/registry"

	"github.com/gogo/protobuf/jsonpb"
)

var (
	DefaultTTL       = time.Hour
	DefaultMarshaler = &jsonpb.Marshaler{}
)

type Base struct {
	Framework dl.FrameworkManifest
}

func toPath(s string) string {
	return strings.Replace(s, ":", "/", -1)
}

func (b *Base) PublishInRegistery(prefix string) error {
	marshaler := DefaultMarshaler
	rgs, err := registry.New()
	if err != nil {
		return err
	}
	defer rgs.Close()

	framework := b.Framework
	cn, err := framework.CanonicalName()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	prefix = path.Join(config.App.Name, prefix)
	rgs.Put(prefix, nil, &store.WriteOptions{IsDir: true})
	{
		frameworksKey := path.Join(prefix, "frameworks")
		// lock, err := rgs.NewLock(frameworksKey, &store.LockOptions{TTL: 2 * time.Second})
		// if err != nil {
		// 	return errors.Wrapf(err, "cannot get lock for key %v", frameworksKey)
		// }
		// pp.Println("trying lock....")
		// _, err = lock.Lock(nil)
		// if err != nil {
		// 	return errors.Wrapf(err, "cannot lock key %v", frameworksKey)
		// }
		// pp.Println("locking....")
		wg.Add(1)
		go func() {
			// defer lock.Unlock()
			defer wg.Done()
			kv, err := rgs.Get(frameworksKey)
			if err != nil {
				if ok, e := rgs.Exists(frameworksKey); e == nil && ok {
					log.WithError(err).Errorf("cannot get value for key %v", frameworksKey)
					return
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
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		key := path.Join(prefix, toPath(cn))
		rgs.Put(key, nil, &store.WriteOptions{TTL: DefaultTTL, IsDir: true})

		key = path.Join(key, "info")
		bts := new(bytes.Buffer)
		err := marshaler.Marshal(bts, &framework)
		if err != nil {
			return
		}
		if err := rgs.Put(key, bts.Bytes(), &store.WriteOptions{TTL: DefaultTTL, IsDir: false}); err != nil {
			return
		}
	}()

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
			bts := new(bytes.Buffer)
			err = marshaler.Marshal(bts, &model)
			// bts, err := model.Marshal()
			if err != nil {
				pp.Println(err)
				return
			}
			key := path.Join(prefix, toPath(mn), "info")
			rgs.Put(key, bts.Bytes(), &store.WriteOptions{TTL: DefaultTTL, IsDir: false})
		}(model)
	}

	wg.Wait()

	return nil
}
