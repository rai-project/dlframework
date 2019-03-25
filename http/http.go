package http

import (
        "io/ioutil"
        "bytes"
        // "io"
        "encoding/json"
        "fmt"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rai-project/dlframework/httpapi/restapi/operations"
	"github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/authentication"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/predict"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/registry"

        // "github.com/justinas/nosurf"
        "github.com/volatiletech/authboss"
        auth "github.com/volatiletech/authboss/auth"
        register "github.com/volatiletech/authboss/register"
        "github.com/volatiletech/authboss/remember"
        "github.com/k0kubun/pp"
)

func ConfigureAPI(api *operations.DlframeworkAPI) http.Handler {
	api.ServeError = ServeError
	api.Logger = log.Debugf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.BasicAuthAuth = func(user string, pass string) (*models.User, error) {
            pp.Println(pass)
		return &models.User{
                    Username: "as29",
                }, nil
	}

        api.AuthenticationLoginHandler = authentication.LoginHandlerFunc(
                func(params authentication.LoginParams, principal *models.User) middleware.Responder {
                        return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer){
                                a := &auth.Auth{ab}
                                req := params.HTTPRequest
                                pp.Println(authboss.GetSession(req, authboss.SessionKey))
                                pp.Println(principal)
                                // b, err := ioutil.ReadAll(req.Body)
                                requestByte, _ := json.Marshal(params.Body)
                                fmt.Println(string(requestByte))
                                req.Body = ioutil.NopCloser(bytes.NewReader(requestByte))

                                pp.Println("Login")
                                a.LoginPost(rw, req)
                                // pp.Println(req.Context().Value(authboss.CTXKeyUser))
                                // u = ab.CurrentUser(req)
                                // fmt.Println(u.GetPID())
                        })
                })
        // api.AuthenticationLoginHandler = http.StripPrefix("/auth", ab.Config.Core.Router)
	// api.AuthenticationSignupHandler = authentication.SignupHandlerFunc(SignupHandler)
        api.AuthenticationSignupHandler = authentication.SignupHandlerFunc(
                func(params authentication.SignupParams) middleware.Responder {
                        return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer){
                                r := &register.Register{ab}
                                req := params.HTTPRequest
                                requestByte, _ := json.Marshal(params.Body)
                                fmt.Println(string(requestByte))
                                req.Body = ioutil.NopCloser(bytes.NewReader(requestByte))
                                r.Post(rw, req)
                        })
                })

	api.RegistryFrameworkAgentsHandler = registry.FrameworkAgentsHandlerFunc(RegistryFrameworkAgentsHandler)
	api.RegistryFrameworkManifestsHandler = registry.FrameworkManifestsHandlerFunc(RegistryFrameworkManifestsHandler)
	api.RegistryModelAgentsHandler = registry.ModelAgentsHandlerFunc(RegistryModelAgentsHandler)
	api.RegistryModelManifestsHandler = registry.ModelManifestsHandlerFunc(RegistryModelManifestsHandler)

	predictHandler := &PredictHandler{}
	api.PredictOpenHandler = predict.OpenHandlerFunc(predictHandler.Open)
	api.PredictCloseHandler = predict.CloseHandlerFunc(predictHandler.Close)
	api.PredictResetHandler = predict.ResetHandlerFunc(predictHandler.Reset)
	api.PredictImagesHandler = predict.ImagesHandlerFunc(predictHandler.Images)
	api.PredictUrlsHandler = predict.UrlsHandlerFunc(predictHandler.URLs)
	api.PredictDatasetHandler = predict.DatasetHandlerFunc(predictHandler.Dataset)

	api.ServerShutdown = func() {}
        setupAuthboss()

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
        // handler = nosurfing(handler)
        handler = ab.LoadClientStateMiddleware(handler)
        handler = remember.Middleware(ab)(handler)
        handler = dataInjector(handler)
	return handler
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}
