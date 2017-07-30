package http

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/registry"
)

func RegistryFrameworkAgentsHandler(params registry.FrameworkAgentsParams) middleware.Responder {
	return middleware.NotImplemented("operation registry.FrameworkAgents has not yet been implemented")
}
func RegistryFrameworkManifestsHandler(params registry.FrameworkManifestsParams) middleware.Responder {
	return middleware.NotImplemented("operation registry.FrameworkManifests has not yet been implemented")
}

func RegistryModelAgentsHandler(params registry.ModelAgentsParams) middleware.Responder {
	return middleware.NotImplemented("operation registry.ModelAgents has not yet been implemented")
}
func RegistryModelManifestsHandler(params registry.ModelManifestsParams) middleware.Responder {
	return middleware.NotImplemented("operation registry.ModelManifests has not yet been implemented")
}
