package dlframework

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/syncmap"
)

var frameworkRegistry = syncmap.Map{}

func (f FrameworkManifest) CannonicalName() (string, error) {
	frameworkName := f.GetName()
	if frameworkName == "" {
		return "", errors.New("framework name must not be empty")
	}
	frameworkVersion := f.GetVersion()
	if frameworkVersion == "" {
		frameworkVersion = "latest"
	}
	return strings.ToLower(frameworkName + ":" + frameworkVersion), nil
}

func (f FrameworkManifest) Register() error {
	name, err := f.CannonicalName()
	if err != nil {
		return err
	}
	return f.RegisterNamed(name)
}

func (f FrameworkManifest) RegisterNamed(s string) error {
	name := strings.ToLower(s)
	if _, ok := frameworkRegistry.Load(name); ok {
		return errors.Errorf("the %s framework has already been registered", s)
	}
	frameworkRegistry.Store(s, f)
	return nil
}

func RegisteredFrameworkNames() []string {
	return syncMapKeys(frameworkRegistry)
}
