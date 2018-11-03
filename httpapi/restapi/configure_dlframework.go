// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/rai-project/dlframework/httpapi/restapi/operations"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/authentication"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/predict"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/registry"
)

//go:generate swagger generate server --target ../httpapi --name Dlframework --spec ../dlframework.swagger.json

func configureFlags(api *operations.DlframeworkAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.DlframeworkAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.PredictCloseHandler = predict.CloseHandlerFunc(func(params predict.CloseParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.Close has not yet been implemented")
	})
	api.PredictDatasetHandler = predict.DatasetHandlerFunc(func(params predict.DatasetParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.Dataset has not yet been implemented")
	})
	api.PredictDatasetStreamHandler = predict.DatasetStreamHandlerFunc(func(params predict.DatasetStreamParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.DatasetStream has not yet been implemented")
	})
	api.RegistryFrameworkAgentsHandler = registry.FrameworkAgentsHandlerFunc(func(params registry.FrameworkAgentsParams) middleware.Responder {
		return middleware.NotImplemented("operation registry.FrameworkAgents has not yet been implemented")
	})
	api.RegistryFrameworkManifestsHandler = registry.FrameworkManifestsHandlerFunc(func(params registry.FrameworkManifestsParams) middleware.Responder {
		return middleware.NotImplemented("operation registry.FrameworkManifests has not yet been implemented")
	})
	api.PredictImagesHandler = predict.ImagesHandlerFunc(func(params predict.ImagesParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.Images has not yet been implemented")
	})
	api.PredictImagesStreamHandler = predict.ImagesStreamHandlerFunc(func(params predict.ImagesStreamParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.ImagesStream has not yet been implemented")
	})
	api.AuthenticationLoginHandler = authentication.LoginHandlerFunc(func(params authentication.LoginParams) middleware.Responder {
		return middleware.NotImplemented("operation authentication.Login has not yet been implemented")
	})
	api.RegistryModelAgentsHandler = registry.ModelAgentsHandlerFunc(func(params registry.ModelAgentsParams) middleware.Responder {
		return middleware.NotImplemented("operation registry.ModelAgents has not yet been implemented")
	})
	api.RegistryModelManifestsHandler = registry.ModelManifestsHandlerFunc(func(params registry.ModelManifestsParams) middleware.Responder {
		return middleware.NotImplemented("operation registry.ModelManifests has not yet been implemented")
	})
	api.PredictOpenHandler = predict.OpenHandlerFunc(func(params predict.OpenParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.Open has not yet been implemented")
	})
	api.PredictResetHandler = predict.ResetHandlerFunc(func(params predict.ResetParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.Reset has not yet been implemented")
	})
	api.AuthenticationSignupHandler = authentication.SignupHandlerFunc(func(params authentication.SignupParams) middleware.Responder {
		return middleware.NotImplemented("operation authentication.Signup has not yet been implemented")
	})
	api.PredictUrlsHandler = predict.UrlsHandlerFunc(func(params predict.UrlsParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.Urls has not yet been implemented")
	})
	api.PredictUrlsStreamHandler = predict.UrlsStreamHandlerFunc(func(params predict.UrlsStreamParams) middleware.Responder {
		return middleware.NotImplemented("operation predict.UrlsStream has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
