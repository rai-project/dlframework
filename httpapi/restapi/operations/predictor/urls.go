// Code generated by go-swagger; DO NOT EDIT.

package predictor

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UrlsHandlerFunc turns a function with the right signature into a urls handler
type UrlsHandlerFunc func(UrlsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UrlsHandlerFunc) Handle(params UrlsParams) middleware.Responder {
	return fn(params)
}

// UrlsHandler interface for that can handle valid urls params
type UrlsHandler interface {
	Handle(UrlsParams) middleware.Responder
}

// NewUrls creates a new http.Handler for the urls operation
func NewUrls(ctx *middleware.Context, handler UrlsHandler) *Urls {
	return &Urls{Context: ctx, Handler: handler}
}

/*Urls swagger:route POST /v1/predict/urls Predictor urls

Urls urls API

*/
type Urls struct {
	Context *middleware.Context
	Handler UrlsHandler
}

func (o *Urls) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUrlsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
