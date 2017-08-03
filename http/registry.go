package http

import (
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/registry"
)

func getParam(val *string, defaultValue string) string {
	if val == nil || *val == "" {
		return defaultValue
	}
	return *val
}

func RegistryFrameworkManifestsHandler(params registry.FrameworkManifestsParams) middleware.Responder {

	sort := func(manifests []*webmodels.DlframeworkFrameworkManifest) []*webmodels.DlframeworkFrameworkManifest {
		sort.Slice(manifests, func(i int, j int) bool {
			a := manifests[i]
			b := manifests[j]
			if a.Name != b.Name {
				return a.Name < b.Name
			}
			if a.Version != b.Version {
				aVersion, _ := semver.NewVersion(a.Version)
				bVersion, _ := semver.NewVersion(b.Version)
				return aVersion.LessThan(bVersion)
			}
			return false
		})
		return manifests
	}

	manifests, err := frameworks.manifests()
	if err != nil {
		return NewError("FrameworkManifests", err)
	}

	if len(manifests) == 0 {
		return NewError("FrameworkManifests",
			errors.New("no frameworks found"),
		)
	}

	frameworkName := strings.ToLower(getParam(params.FrameworkName, "*"))
	frameworkVersion := strings.ToLower(getParam(params.FrameworkVersion, "*"))

	manifests, err = frameworks.filterManifests(manifests, frameworkName, frameworkVersion)
	if err != nil {
		return NewError("FrameworkManifests", err)
	}

	return registry.NewFrameworkManifestsOK().
		WithPayload(&webmodels.DlframeworkFrameworkManifestsResponse{
			Manifests: sort(manifests),
		})
}

func RegistryModelManifestsHandler(params registry.ModelManifestsParams) middleware.Responder {

	sort := func(manifests []*webmodels.DlframeworkModelManifest) []*webmodels.DlframeworkModelManifest {
		sort.Slice(manifests, func(i int, j int) bool {
			a := manifests[i]
			b := manifests[j]
			if a.Framework.Name != b.Framework.Name {
				return a.Framework.Name < b.Framework.Name
			}
			if a.Name != b.Name {
				return a.Name < b.Name
			}
			if a.Version != b.Version {
				aVersion, _ := semver.NewVersion(a.Version)
				bVersion, _ := semver.NewVersion(b.Version)
				return aVersion.LessThan(bVersion)
			}
			return false
		})
		return manifests
	}

	frameworkName := strings.ToLower(getParam(params.FrameworkName, "*"))
	frameworkVersion := strings.ToLower(getParam(params.FrameworkVersion, "*"))

	manifests, err := models.manifests(frameworkName, frameworkVersion)
	if err != nil {
		return NewError("ModelManifests", err)
	}

	if len(manifests) == 0 {
		return NewError("ModelManifests",
			errors.Errorf("no models found for the framework %s:%s", frameworkName, frameworkVersion),
		)
	}

	modelName := strings.ToLower(getParam(params.ModelName, "*"))
	modelVersion := strings.ToLower(getParam(params.ModelVersion, "*"))

	manifests, err = models.filterManifests(manifests, modelName, modelVersion)
	if err != nil {
		return NewError("ModelManifests", err)
	}

	return registry.NewModelManifestsOK().
		WithPayload(&webmodels.DlframeworkModelManifestsResponse{
			Manifests: sort(manifests),
		})
}

func RegistryFrameworkAgentsHandler(params registry.FrameworkAgentsParams) middleware.Responder {
	frameworkName := strings.ToLower(getParam(params.FrameworkName, "*"))
	frameworkVersion := strings.ToLower(getParam(params.FrameworkVersion, "*"))
	modelName := "*"
	modelVersion := "*"

	agents, err := models.agents(frameworkName, frameworkVersion, modelName, modelVersion)
	if err != nil {
		return NewError("ModelAgents", err)
	}

	return registry.NewFrameworkAgentsOK().
		WithPayload(&webmodels.DlframeworkAgents{
			Agents: agents,
		})
}

func RegistryModelAgentsHandler(params registry.ModelAgentsParams) middleware.Responder {

	frameworkName := strings.ToLower(getParam(params.FrameworkName, "*"))
	frameworkVersion := strings.ToLower(getParam(params.FrameworkVersion, "*"))
	modelName := strings.ToLower(getParam(params.ModelName, "*"))
	modelVersion := strings.ToLower(getParam(params.ModelVersion, "*"))

	agents, err := models.agents(frameworkName, frameworkVersion, modelName, modelVersion)
	if err != nil {
		return NewError("ModelAgents", err)
	}

	return registry.NewModelAgentsOK().
		WithPayload(&webmodels.DlframeworkAgents{
			Agents: agents,
		})
}
