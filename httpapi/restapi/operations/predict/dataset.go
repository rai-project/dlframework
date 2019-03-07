// Code generated by go-swagger; DO NOT EDIT.

package predict

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DatasetHandlerFunc turns a function with the right signature into a dataset handler
type DatasetHandlerFunc func(DatasetParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DatasetHandlerFunc) Handle(params DatasetParams) middleware.Responder {
	return fn(params)
}

// DatasetHandler interface for that can handle valid dataset params
type DatasetHandler interface {
	Handle(DatasetParams) middleware.Responder
}

// NewDataset creates a new http.Handler for the dataset operation
func NewDataset(ctx *middleware.Context, handler DatasetHandler) *Dataset {
	return &Dataset{Context: ctx, Handler: handler}
}

/*Dataset swagger:route POST /predict/dataset Predict dataset

Dataset method receives a single dataset and runs
the predictor on all elements of the dataset.
The result is a prediction feature list.

*/
type Dataset struct {
	Context *middleware.Context
	Handler DatasetHandler
}

func (o *Dataset) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDatasetParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
