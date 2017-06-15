package registry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetFrameworkManifestParams creates a new GetFrameworkManifestParams object
// with the default values initialized.
func NewGetFrameworkManifestParams() GetFrameworkManifestParams {
	var ()
	return GetFrameworkManifestParams{}
}

// GetFrameworkManifestParams contains all the bound params for the get framework manifest operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetFrameworkManifest
type GetFrameworkManifestParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*
	  Required: true
	  In: path
	*/
	FrameworkName string
	/*
	  In: query
	*/
	FrameworkVersion *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetFrameworkManifestParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rFrameworkName, rhkFrameworkName, _ := route.Params.GetOK("framework_name")
	if err := o.bindFrameworkName(rFrameworkName, rhkFrameworkName, route.Formats); err != nil {
		res = append(res, err)
	}

	qFrameworkVersion, qhkFrameworkVersion, _ := qs.GetOK("framework_version")
	if err := o.bindFrameworkVersion(qFrameworkVersion, qhkFrameworkVersion, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetFrameworkManifestParams) bindFrameworkName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.FrameworkName = raw

	return nil
}

func (o *GetFrameworkManifestParams) bindFrameworkVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.FrameworkVersion = &raw

	return nil
}
