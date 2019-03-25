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

	models "github.com/rai-project/dlframework/httpapi/models"
)

//go:generate swagger generate server --target ../../httpapi --name Dlframework --spec ../../dlframework.swagger.json --principal models.User

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

	// Applies when the Authorization header is set with the Basic scheme
	api.BasicAuthAuth = func(user string, pass string) (*models.User, error) {
		return nil, errors.NotImplemented("basic auth  (basicAuth) has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	if api.PredictCloseHandler == nil {
		api.PredictCloseHandler = predict.CloseHandlerFunc(func(params predict.CloseParams) middleware.Responder {
			return middleware.NotImplemented("operation predict.Close has not yet been implemented")
		})
	}
	if api.PredictDatasetHandler == nil {
		api.PredictDatasetHandler = predict.DatasetHandlerFunc(func(params predict.DatasetParams) middleware.Responder {
			return middleware.NotImplemented("operation predict.Dataset has not yet been implemented")
		})
	}
	if api.RegistryFrameworkAgentsHandler == nil {
		api.RegistryFrameworkAgentsHandler = registry.FrameworkAgentsHandlerFunc(func(params registry.FrameworkAgentsParams) middleware.Responder {
			return middleware.NotImplemented("operation registry.FrameworkAgents has not yet been implemented")
		})
	}
	if api.RegistryFrameworkManifestsHandler == nil {
		api.RegistryFrameworkManifestsHandler = registry.FrameworkManifestsHandlerFunc(func(params registry.FrameworkManifestsParams) middleware.Responder {
			return middleware.NotImplemented("operation registry.FrameworkManifests has not yet been implemented")
		})
	}
	if api.PredictImagesHandler == nil {
		api.PredictImagesHandler = predict.ImagesHandlerFunc(func(params predict.ImagesParams) middleware.Responder {
			return middleware.NotImplemented("operation predict.Images has not yet been implemented")
		})
	}
	if api.AuthenticationLoginHandler == nil {
		api.AuthenticationLoginHandler = authentication.LoginHandlerFunc(func(params authentication.LoginParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation authentication.Login has not yet been implemented")
		})
	}
	if api.AuthenticationLogoutHandler == nil {
		api.AuthenticationLogoutHandler = authentication.LogoutHandlerFunc(func(params authentication.LogoutParams) middleware.Responder {
			return middleware.NotImplemented("operation authentication.Logout has not yet been implemented")
		})
	}
	if api.RegistryModelAgentsHandler == nil {
		api.RegistryModelAgentsHandler = registry.ModelAgentsHandlerFunc(func(params registry.ModelAgentsParams) middleware.Responder {
			return middleware.NotImplemented("operation registry.ModelAgents has not yet been implemented")
		})
	}
	if api.RegistryModelManifestsHandler == nil {
		api.RegistryModelManifestsHandler = registry.ModelManifestsHandlerFunc(func(params registry.ModelManifestsParams) middleware.Responder {
			return middleware.NotImplemented("operation registry.ModelManifests has not yet been implemented")
		})
	}
	if api.PredictOpenHandler == nil {
		api.PredictOpenHandler = predict.OpenHandlerFunc(func(params predict.OpenParams) middleware.Responder {
			return middleware.NotImplemented("operation predict.Open has not yet been implemented")
		})
	}
	if api.PredictResetHandler == nil {
		api.PredictResetHandler = predict.ResetHandlerFunc(func(params predict.ResetParams) middleware.Responder {
			return middleware.NotImplemented("operation predict.Reset has not yet been implemented")
		})
	}
	if api.AuthenticationSignupHandler == nil {
		api.AuthenticationSignupHandler = authentication.SignupHandlerFunc(func(params authentication.SignupParams) middleware.Responder {
			return middleware.NotImplemented("operation authentication.Signup has not yet been implemented")
		})
	}
	if api.PredictUrlsHandler == nil {
		api.PredictUrlsHandler = predict.UrlsHandlerFunc(func(params predict.UrlsParams) middleware.Responder {
			return middleware.NotImplemented("operation predict.Urls has not yet been implemented")
		})
	}
	if api.AuthenticationUpdateHandler == nil {
		api.AuthenticationUpdateHandler = authentication.UpdateHandlerFunc(func(params authentication.UpdateParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation authentication.Update has not yet been implemented")
		})
	}
	if api.AuthenticationUserInfoHandler == nil {
		api.AuthenticationUserInfoHandler = authentication.UserInfoHandlerFunc(func(params authentication.UserInfoParams) middleware.Responder {
			return middleware.NotImplemented("operation authentication.UserInfo has not yet been implemented")
		})
	}

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
