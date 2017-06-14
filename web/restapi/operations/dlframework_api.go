package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	loads "github.com/go-openapi/loads"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	spec "github.com/go-openapi/spec"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/rai-project/dlframework/web/restapi/operations/carml"
	"github.com/rai-project/dlframework/web/restapi/operations/predictor"
)

// NewDlframeworkAPI creates a new Dlframework instance
func NewDlframeworkAPI(spec *loads.Document) *DlframeworkAPI {
	return &DlframeworkAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,
		JSONConsumer:        runtime.JSONConsumer(),
		JSONProducer:        runtime.JSONProducer(),
		CarmlGetFrameworkManifestHandler: carml.GetFrameworkManifestHandlerFunc(func(params carml.GetFrameworkManifestParams) middleware.Responder {
			return middleware.NotImplemented("operation CarmlGetFrameworkManifest has not yet been implemented")
		}),
		CarmlGetFrameworkManifestsHandler: carml.GetFrameworkManifestsHandlerFunc(func(params carml.GetFrameworkManifestsParams) middleware.Responder {
			return middleware.NotImplemented("operation CarmlGetFrameworkManifests has not yet been implemented")
		}),
		CarmlGetFrameworkModelManifestHandler: carml.GetFrameworkModelManifestHandlerFunc(func(params carml.GetFrameworkModelManifestParams) middleware.Responder {
			return middleware.NotImplemented("operation CarmlGetFrameworkModelManifest has not yet been implemented")
		}),
		CarmlGetFrameworkModelsHandler: carml.GetFrameworkModelsHandlerFunc(func(params carml.GetFrameworkModelsParams) middleware.Responder {
			return middleware.NotImplemented("operation CarmlGetFrameworkModels has not yet been implemented")
		}),
		CarmlGetModelManifestHandler: carml.GetModelManifestHandlerFunc(func(params carml.GetModelManifestParams) middleware.Responder {
			return middleware.NotImplemented("operation CarmlGetModelManifest has not yet been implemented")
		}),
		CarmlGetModelManifestsHandler: carml.GetModelManifestsHandlerFunc(func(params carml.GetModelManifestsParams) middleware.Responder {
			return middleware.NotImplemented("operation CarmlGetModelManifests has not yet been implemented")
		}),
		PredictorPredictHandler: predictor.PredictHandlerFunc(func(params predictor.PredictParams) middleware.Responder {
			return middleware.NotImplemented("operation PredictorPredict has not yet been implemented")
		}),
	}
}

/*DlframeworkAPI TODO... fillme. */
type DlframeworkAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for a "application/json" mime type
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for a "application/json" mime type
	JSONProducer runtime.Producer

	// CarmlGetFrameworkManifestHandler sets the operation handler for the get framework manifest operation
	CarmlGetFrameworkManifestHandler carml.GetFrameworkManifestHandler
	// CarmlGetFrameworkManifestsHandler sets the operation handler for the get framework manifests operation
	CarmlGetFrameworkManifestsHandler carml.GetFrameworkManifestsHandler
	// CarmlGetFrameworkModelManifestHandler sets the operation handler for the get framework model manifest operation
	CarmlGetFrameworkModelManifestHandler carml.GetFrameworkModelManifestHandler
	// CarmlGetFrameworkModelsHandler sets the operation handler for the get framework models operation
	CarmlGetFrameworkModelsHandler carml.GetFrameworkModelsHandler
	// CarmlGetModelManifestHandler sets the operation handler for the get model manifest operation
	CarmlGetModelManifestHandler carml.GetModelManifestHandler
	// CarmlGetModelManifestsHandler sets the operation handler for the get model manifests operation
	CarmlGetModelManifestsHandler carml.GetModelManifestsHandler
	// PredictorPredictHandler sets the operation handler for the predict operation
	PredictorPredictHandler predictor.PredictHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// SetDefaultProduces sets the default produces media type
func (o *DlframeworkAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *DlframeworkAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *DlframeworkAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *DlframeworkAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *DlframeworkAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *DlframeworkAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *DlframeworkAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the DlframeworkAPI
func (o *DlframeworkAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.CarmlGetFrameworkManifestHandler == nil {
		unregistered = append(unregistered, "carml.GetFrameworkManifestHandler")
	}

	if o.CarmlGetFrameworkManifestsHandler == nil {
		unregistered = append(unregistered, "carml.GetFrameworkManifestsHandler")
	}

	if o.CarmlGetFrameworkModelManifestHandler == nil {
		unregistered = append(unregistered, "carml.GetFrameworkModelManifestHandler")
	}

	if o.CarmlGetFrameworkModelsHandler == nil {
		unregistered = append(unregistered, "carml.GetFrameworkModelsHandler")
	}

	if o.CarmlGetModelManifestHandler == nil {
		unregistered = append(unregistered, "carml.GetModelManifestHandler")
	}

	if o.CarmlGetModelManifestsHandler == nil {
		unregistered = append(unregistered, "carml.GetModelManifestsHandler")
	}

	if o.PredictorPredictHandler == nil {
		unregistered = append(unregistered, "predictor.PredictHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *DlframeworkAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *DlframeworkAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {

	return nil

}

// ConsumersFor gets the consumers for the specified media types
func (o *DlframeworkAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {

	result := make(map[string]runtime.Consumer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONConsumer

		}
	}
	return result

}

// ProducersFor gets the producers for the specified media types
func (o *DlframeworkAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {

	result := make(map[string]runtime.Producer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONProducer

		}
	}
	return result

}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *DlframeworkAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the dlframework API
func (o *DlframeworkAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *DlframeworkAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened

	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/v1/framework/{framework_name}/{framework_version}/info"] = carml.NewGetFrameworkManifest(o.context, o.CarmlGetFrameworkManifestHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/v1/frameworks"] = carml.NewGetFrameworkManifests(o.context, o.CarmlGetFrameworkManifestsHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/v1/framework/{framework_name}/{framework_version}/model/{model_name}/{model_version}/info"] = carml.NewGetFrameworkModelManifest(o.context, o.CarmlGetFrameworkModelManifestHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/v1/framework/{framework_name}/{framework_version}/models"] = carml.NewGetFrameworkModels(o.context, o.CarmlGetFrameworkModelsHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/v1/model/{model_name}/{model_version}/info"] = carml.NewGetModelManifest(o.context, o.CarmlGetModelManifestHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/v1/models"] = carml.NewGetModelManifests(o.context, o.CarmlGetModelManifestsHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/v1/{framework_name}/{framework_version}/{model_name}/{model_version}/predict"] = predictor.NewPredict(o.context, o.PredictorPredictHandler)

}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *DlframeworkAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middelware as you see fit
func (o *DlframeworkAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}
