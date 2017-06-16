package dlframework

import (
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"golang.org/x/sync/syncmap"
)

var frameworkRegistry = syncmap.Map{}

func (f FrameworkManifest) MustCanonicalName() string {
	s, err := f.CanonicalName()
	if err != nil {
		log.WithField("model_name", f.GetName()).Fatal("unable to get framework canonical name")
		return ""
	}
	return s
}

func (f FrameworkManifest) CanonicalName() (string, error) {
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
	name, err := f.CanonicalName()
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

func Frameworks() ([]FrameworkManifest, error) {
	names := RegisteredFrameworkNames()
	fws := make([]FrameworkManifest, len(names))
	for ii, name := range names {
		f, ok := frameworkRegistry.Load(name)
		if !ok {
			return nil, errors.Errorf("framework %s was not found", name)
		}
		fw, ok := f.(FrameworkManifest)
		if !ok {
			return nil, errors.Errorf("framework %s was found but not of type FrameworkManifest", name)
		}
		fws[ii] = fw
	}
	return fws, nil
}

func (f FrameworkManifest) Models() []ModelManifest { // todo: this is not optimal
	models := []ModelManifest{}
	names := RegisteredModelNames()
	for _, name := range names {
		m, err := f.FindModel(name)
		if err != nil {
			continue
		}
		models = append(models, m)
	}
	return models
}

func (f FrameworkManifest) FindModel(name string) (ModelManifest, error) {
	var model *ModelManifest
	frameworkVersionString := f.GetVersion()
	if frameworkVersionString == "" {
		return ModelManifest{}, errors.Errorf("expecting a framework version for framework = %v", f.GetName())
	}
	frameworkVersion, err := semver.NewVersion(frameworkVersionString)
	if err != nil {
		return ModelManifest{}, errors.Wrapf(err, "unable to parse version information for framework = %v", f.GetName())
	}
	frameworkCanonicalName, err := f.CanonicalName()
	if err != nil {
		return ModelManifest{}, err
	}
	name = strings.ToLower(name)
	modelRegistry.Range(func(key0 interface{}, value interface{}) bool {
		key, ok := key0.(string)
		if !ok {
			return true
		}
		key = strings.TrimPrefix(key, frameworkCanonicalName+"/")
		if key != name {
			return true
		}
		m, ok := value.(ModelManifest)
		if !ok {
			return true
		}
		if m.GetFramework().GetVersion() != "latest" {
			cs, err := m.FrameworkConstraint()
			if err != nil {
				return true
			}
			ok := cs.Check(frameworkVersion)
			if !ok {
				return true
			}
		}
		model = &m
		return false
	})
	if model == nil {
		return ModelManifest{}, errors.Errorf("model %s for framework %s not found in registry", name, f.GetName())
	}
	return *model, nil
}

func FindFramework(name string) (*FrameworkManifest, error) {
	var framework *FrameworkManifest
	modelRegistry.Range(func(key0 interface{}, value interface{}) bool {
		key, ok := key0.(string)
		if !ok {
			return true
		}
		if key != name {
			return true
		}
		f, ok := value.(FrameworkManifest)
		if !ok {
			return true
		}
		framework = &f
		return false
	})
	if framework == nil {
		return nil, errors.Errorf("framework %s not found in registry", name)
	}
	return framework, nil
}
