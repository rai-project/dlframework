package carml

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/rai-project/dlframework/web/models"
)

// NewPredictParams creates a new PredictParams object
// with the default values initialized.
func NewPredictParams() *PredictParams {
	var ()
	return &PredictParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPredictParamsWithTimeout creates a new PredictParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPredictParamsWithTimeout(timeout time.Duration) *PredictParams {
	var ()
	return &PredictParams{

		timeout: timeout,
	}
}

// NewPredictParamsWithContext creates a new PredictParams object
// with the default values initialized, and the ability to set a context for a request
func NewPredictParamsWithContext(ctx context.Context) *PredictParams {
	var ()
	return &PredictParams{

		Context: ctx,
	}
}

// NewPredictParamsWithHTTPClient creates a new PredictParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPredictParamsWithHTTPClient(client *http.Client) *PredictParams {
	var ()
	return &PredictParams{
		HTTPClient: client,
	}
}

/*PredictParams contains all the parameters to send to the API endpoint
for the predict operation typically these are written to a http.Request
*/
type PredictParams struct {

	/*Body*/
	Body *models.DlframeworkPredictRequest
	/*FrameworkName*/
	FrameworkName string
	/*FrameworkVersion*/
	FrameworkVersion string
	/*ModelName*/
	ModelName string
	/*ModelVersion*/
	ModelVersion string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the predict params
func (o *PredictParams) WithTimeout(timeout time.Duration) *PredictParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the predict params
func (o *PredictParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the predict params
func (o *PredictParams) WithContext(ctx context.Context) *PredictParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the predict params
func (o *PredictParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the predict params
func (o *PredictParams) WithHTTPClient(client *http.Client) *PredictParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the predict params
func (o *PredictParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the predict params
func (o *PredictParams) WithBody(body *models.DlframeworkPredictRequest) *PredictParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the predict params
func (o *PredictParams) SetBody(body *models.DlframeworkPredictRequest) {
	o.Body = body
}

// WithFrameworkName adds the frameworkName to the predict params
func (o *PredictParams) WithFrameworkName(frameworkName string) *PredictParams {
	o.SetFrameworkName(frameworkName)
	return o
}

// SetFrameworkName adds the frameworkName to the predict params
func (o *PredictParams) SetFrameworkName(frameworkName string) {
	o.FrameworkName = frameworkName
}

// WithFrameworkVersion adds the frameworkVersion to the predict params
func (o *PredictParams) WithFrameworkVersion(frameworkVersion string) *PredictParams {
	o.SetFrameworkVersion(frameworkVersion)
	return o
}

// SetFrameworkVersion adds the frameworkVersion to the predict params
func (o *PredictParams) SetFrameworkVersion(frameworkVersion string) {
	o.FrameworkVersion = frameworkVersion
}

// WithModelName adds the modelName to the predict params
func (o *PredictParams) WithModelName(modelName string) *PredictParams {
	o.SetModelName(modelName)
	return o
}

// SetModelName adds the modelName to the predict params
func (o *PredictParams) SetModelName(modelName string) {
	o.ModelName = modelName
}

// WithModelVersion adds the modelVersion to the predict params
func (o *PredictParams) WithModelVersion(modelVersion string) *PredictParams {
	o.SetModelVersion(modelVersion)
	return o
}

// SetModelVersion adds the modelVersion to the predict params
func (o *PredictParams) SetModelVersion(modelVersion string) {
	o.ModelVersion = modelVersion
}

// WriteToRequest writes these params to a swagger request
func (o *PredictParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body == nil {
		o.Body = new(models.DlframeworkPredictRequest)
	}

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	// path param framework_name
	if err := r.SetPathParam("framework_name", o.FrameworkName); err != nil {
		return err
	}

	// path param framework_version
	if err := r.SetPathParam("framework_version", o.FrameworkVersion); err != nil {
		return err
	}

	// path param model_name
	if err := r.SetPathParam("model_name", o.ModelName); err != nil {
		return err
	}

	// path param model_version
	if err := r.SetPathParam("model_version", o.ModelVersion); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
