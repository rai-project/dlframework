package carml

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetFrameworkModelManifestHandlerFunc turns a function with the right signature into a get framework model manifest handler
type GetFrameworkModelManifestHandlerFunc func(GetFrameworkModelManifestParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetFrameworkModelManifestHandlerFunc) Handle(params GetFrameworkModelManifestParams) middleware.Responder {
	return fn(params)
}

// GetFrameworkModelManifestHandler interface for that can handle valid get framework model manifest params
type GetFrameworkModelManifestHandler interface {
	Handle(GetFrameworkModelManifestParams) middleware.Responder
}

// NewGetFrameworkModelManifest creates a new http.Handler for the get framework model manifest operation
func NewGetFrameworkModelManifest(ctx *middleware.Context, handler GetFrameworkModelManifestHandler) *GetFrameworkModelManifest {
	return &GetFrameworkModelManifest{Context: ctx, Handler: handler}
}

/*GetFrameworkModelManifest swagger:route POST /v1/framework/{framework_name}/{framework_version}/model/{model_name}/{model_version}/info carml getFrameworkModelManifest

GetFrameworkModelManifest get framework model manifest API

*/
type GetFrameworkModelManifest struct {
	Context *middleware.Context
	Handler GetFrameworkModelManifestHandler
}

func (o *GetFrameworkModelManifest) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetFrameworkModelManifestParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
