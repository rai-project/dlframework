package http

import (
	"sort"
	"strings"

	"github.com/k0kubun/pp"

	"github.com/Masterminds/semver"
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/registry"
	"github.com/rai-project/dlframework/registryquery"
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
			aName := strings.ToLower(a.Name)
			bName := strings.ToLower(b.Name)
			if aName != bName {
				return aName < bName
			}
			if a.Version != b.Version {
				aVersion, err := semver.NewVersion(a.Version)
				if err != nil {
					return false
				}
				bVersion, err := semver.NewVersion(b.Version)
				if err != nil {
					return true
				}
				return aVersion.LessThan(bVersion)
			}
			return false
		})
		return manifests
	}

	manifests, err := registryquery.Frameworks.Manifests()
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

	manifests, err = registryquery.Frameworks.FilterManifests(manifests, frameworkName, frameworkVersion)
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
			aFrameworkName := strings.ToLower(a.Framework.Name)
			bFrameworkName := strings.ToLower(b.Framework.Name)
			if aFrameworkName != bFrameworkName {
				return aFrameworkName < bFrameworkName
			}
			aName := strings.ToLower(a.Name)
			bName := strings.ToLower(b.Name)
			if aName != bName {
				return aName < bName
			}
			if a.Version != b.Version {
				aVersion, err := semver.NewVersion(a.Version)
				if err != nil {
					return false
				}
				bVersion, err := semver.NewVersion(b.Version)
				if err != nil {
					return true
				}
				return aVersion.LessThan(bVersion)
			}
			return false
		})
		return manifests
	}

	frameworkName := strings.ToLower(getParam(params.FrameworkName, "*"))
	frameworkVersion := strings.ToLower(getParam(params.FrameworkVersion, "*"))

	manifests, err := registryquery.Models.Manifests(frameworkName, frameworkVersion)
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

	manifests, err = registryquery.Models.FilterManifests(manifests, modelName, modelVersion)
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

	agents, err := registryquery.Models.Agents(frameworkName, frameworkVersion, modelName, modelVersion)
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

	agents, err := registryquery.Models.Agents(frameworkName, frameworkVersion, modelName, modelVersion)
	pp.Println(agents)
	if err != nil {
		return NewError("ModelAgents", err)
	}

	return registry.NewModelAgentsOK().
		WithPayload(&webmodels.DlframeworkAgents{
			Agents: agents,
		})
}
