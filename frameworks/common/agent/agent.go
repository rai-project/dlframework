package agent

import (
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rai-project/config"
	dl "github.com/rai-project/dlframework"
	store "github.com/rai-project/libkv/store"
	"github.com/rai-project/registry"
)

type Base struct {
	Framework dl.FrameworkManifest
}

func toPath(s string) string {
	return strings.Replace(s, ":", "/", -1)
}

func (b *Base) PublishInRegistery(prefix string) error {
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

	prefix = path.Join(config.App.Name, prefix)
	{
		frameworksKey := path.Join(prefix, "frameworks")
		lk, err := rgs.NewLock(frameworksKey, &store.LockOptions{TTL: time.Second})
		if err != nil {
			return errors.Wrapf(err, "cannot get lock for %v", frameworksKey)
		}
		_, err = lk.Lock(nil)
		if err != nil {
			return errors.Wrapf(err, "cannot lock %v", frameworksKey)
		}
		kv, err := rgs.Get(frameworksKey)
		if err != nil {
			return errors.Wrapf(err, "cannot get value for key %v", frameworksKey)
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
			rgs.Put(frameworksKey, []byte(newVal), &store.WriteOptions{IsDir: false})
		}
		lk.Unlock()
	}
	key := path.Join(config.App.Name, prefix, toPath(cn))
	if err := rgs.Put(key, []byte(cn), nil); err != nil {
		return err
	}

	return nil
}
