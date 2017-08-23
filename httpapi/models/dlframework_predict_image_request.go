// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// DlframeworkPredictImageRequest dlframework predict image request
// swagger:model dlframeworkPredictImageRequest
type DlframeworkPredictImageRequest struct {

	// framework name
	FrameworkName string `json:"framework_name,omitempty"`

	// framework version
	FrameworkVersion string `json:"framework_version,omitempty"`

	// image
	Image strfmt.Base64 `json:"image,omitempty"`

	// limit
	Limit int32 `json:"limit,omitempty"`

	// model name
	ModelName string `json:"model_name,omitempty"`

	// model version
	ModelVersion string `json:"model_version,omitempty"`
}

// Validate validates this dlframework predict image request
func (m *DlframeworkPredictImageRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *DlframeworkPredictImageRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DlframeworkPredictImageRequest) UnmarshalBinary(b []byte) error {
	var res DlframeworkPredictImageRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
