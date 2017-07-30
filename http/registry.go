package http

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/registry"
)

func RegistryFrameworkManifestsHandler(params registry.FrameworkManifestsParams) middleware.Responder {

	manifests, err := frameworks.manifests()
	if err != nil {
		return NewError("FrameworkManifests", err)
	}

	frameworkName := "*"
	if params.FrameworkName != nil && *params.FrameworkName != "" {
		frameworkName = *params.FrameworkName
	}
	frameworkName = strings.ToLower(frameworkName)
	frameworkVersion := "*"
	if params.FrameworkVersion != nil && *params.FrameworkVersion != "" {
		frameworkName = *params.FrameworkVersion
	}
	frameworkVersion = strings.ToLower(frameworkVersion)

	manifests, err = frameworks.filter(manifests, frameworkName, frameworkVersion)
	if err != nil {
		return NewError("FrameworkManifests", err)
	}

	return registry.NewFrameworkManifestsOK().
		WithPayload(&webmodels.DlframeworkFrameworkManifestsResponse{
			Manifests: manifests,
		})
}

func RegistryModelManifestsHandler(params registry.ModelManifestsParams) middleware.Responder {

	frameworkName := "*"
	if params.FrameworkName != nil && *params.FrameworkName != "" {
		frameworkName = *params.FrameworkName
	}
	frameworkName = strings.ToLower(frameworkName)
	frameworkVersion := "*"
	if params.FrameworkVersion != nil && *params.FrameworkVersion != "" {
		frameworkName = *params.FrameworkVersion
	}
	frameworkVersion = strings.ToLower(frameworkVersion)

	manifests, err := models.manifests(frameworkName, frameworkVersion)
	if err != nil {
		return NewError("FrameworkManifests", err)
	}

	modelName := "*"
	if params.ModelName != nil && *params.ModelName != "" {
		modelName = *params.ModelName
	}
	modelName = strings.ToLower(modelName)
	modelVersion := "*"
	if params.ModelVersion != nil && *params.ModelVersion != "" {
		modelVersion = *params.ModelVersion
	}
	modelVersion = strings.ToLower(modelVersion)

	manifests, err = models.filter(manifests, modelName, modelVersion)
	if err != nil {
		return NewError("ModelManifests", err)
	}

	return registry.NewModelManifestsOK().
		WithPayload(&webmodels.DlframeworkModelManifestsResponse{
			Manifests: manifests,
		})
}

func RegistryFrameworkAgentsHandler(params registry.FrameworkAgentsParams) middleware.Responder {
	return middleware.NotImplemented("operation registry.FrameworkAgents has not yet been implemented")
}

func RegistryModelAgentsHandler(params registry.ModelAgentsParams) middleware.Responder {
	return middleware.NotImplemented("operation registry.ModelAgents has not yet been implemented")
}
