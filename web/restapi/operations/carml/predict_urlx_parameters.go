package carml

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewPredictUrlxParams creates a new PredictUrlxParams object
// with the default values initialized.
func NewPredictUrlxParams() PredictUrlxParams {
	var ()
	return PredictUrlxParams{}
}

// PredictUrlxParams contains all the bound params for the predict urlx operation
// typically these are obtained from a http.Request
//
// swagger:parameters PredictURLx
type PredictUrlxParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*
	  Required: true
	  In: body
	*/
	Body *string
	/*
	  Required: true
	  In: path
	*/
	FrameworkName string
	/*
	  Required: true
	  In: path
	*/
	ModelName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *PredictUrlxParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body string
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}

		} else {

			if len(res) == 0 {
				o.Body = &body
			}
		}

	} else {
		res = append(res, errors.Required("body", "body"))
	}

	rFrameworkName, rhkFrameworkName, _ := route.Params.GetOK("framework_name")
	if err := o.bindFrameworkName(rFrameworkName, rhkFrameworkName, route.Formats); err != nil {
		res = append(res, err)
	}

	rModelName, rhkModelName, _ := route.Params.GetOK("model_name")
	if err := o.bindModelName(rModelName, rhkModelName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PredictUrlxParams) bindFrameworkName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.FrameworkName = raw

	return nil
}

func (o *PredictUrlxParams) bindModelName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.ModelName = raw

	return nil
}
