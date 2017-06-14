package dlframework

import (
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"golang.org/x/sync/syncmap"
)

var modelRegistry = syncmap.Map{}

func (m *ModelManifest) Validate() error {
	name := m.GetName()
	if name == "" {
		return errors.New("model name cannot be empty")
	}
	version := m.GetVersion()
	if version == "" {
		version = "latest"
	}
	if version != "latest" {
		if _, err := semver.NewVersion(version); err != nil {
			return errors.Wrapf(err,
				"the version %s for the model %s is not in semantic version format",
				m.GetVersion(),
				m.GetName(),
			)
		}
	}
	return nil
}

func (m ModelManifest) CannonicalName() (string, error) {
	if m.GetFramework() == nil {
		return "", errors.Errorf("the model %s does not have a valid framework", m.GetName())
	}
	frameworkName, err := m.GetFramework().CannonicalName()
	if err != nil {
		return "", errors.Wrapf(err, "cannot get cannonical name for the framework %s and model %s in the registry", m.GetFramework().GetName(), m.GetName())
	}
	fm, ok := frameworkRegistry.Load(frameworkName)
	if !ok {
		return "", errors.Wrapf(err, "cannot get frame %s for model %s in the registry", frameworkName, m.GetName())
	}
	if _, ok := fm.(FrameworkManifest); !ok {
		return "", errors.Errorf("invalid framework %s registered for model %s in the registry", frameworkName, m.GetName())
	}
	modelName := m.GetName()
	if modelName == "" {
		return "", errors.New("model name must not be empty")
	}
	modelVersion := m.GetVersion()
	if modelVersion == "" {
		modelVersion = "latest"
	}
	return frameworkName + "/" + modelName + ":" + modelVersion, nil
}

func (m ModelManifest) Register() error {
	n, err := m.CannonicalName()
	if err != nil {
		return err
	}
	return m.RegisterNamed(n)
}

func (m ModelManifest) RegisterNamed(s string) error {
	name := strings.ToLower(s)
	if _, ok := modelRegistry.Load(name); ok {
		return errors.Errorf("the %s model has already been registered", s)
	}
	modelRegistry.Store(s, m)
	return nil
}

func RegisteredModelNames() []string {
	return syncMapKeys(modelRegistry)
}

func Models() ([]ModelManifest, error) {
	names := RegisteredModelNames()
	models := make([]ModelManifest, len(names))
	for ii, name := range names {
		m, ok := modelRegistry.Load(name)
		if !ok {
			return nil, errors.Errorf("model %s was not found", name)
		}
		model, ok := m.(ModelManifest)
		if !ok {
			return nil, errors.Errorf("model %s was found but not of type ModelManifest", name)
		}
		models[ii] = model
	}
	return models, nil
}

func FindModel(name string) (*ModelManifest, error) {
	var model *ModelManifest
	modelRegistry.Range(func(key0 interface{}, value interface{}) bool {
		key, ok := key0.(string)
		if !ok {
			return true
		}
		if key != name {
			return true
		}
		m, ok := value.(ModelManifest)
		if !ok {
			return true
		}
		model = &m
		return true
	})
	if model == nil {
		return nil, errors.Errorf("model %s not found in registery", name)
	}
	return model, nil
}
