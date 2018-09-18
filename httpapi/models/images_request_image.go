// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ImagesRequestImage images request image
// swagger:model ImagesRequestImage
type ImagesRequestImage struct {

	// The image is base64 encoded
	// Format: byte
	Data strfmt.Base64 `json:"data,omitempty"`

	// An id used to identify the output feature: maps to input_id for output
	ID string `json:"id,omitempty"`
}

// Validate validates this images request image
func (m *ImagesRequestImage) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ImagesRequestImage) validateData(formats strfmt.Registry) error {

	if swag.IsZero(m.Data) { // not required
		return nil
	}

	// Format "byte" (base64 string) is already validated when unmarshalled

	return nil
}

// MarshalBinary interface implementation
func (m *ImagesRequestImage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ImagesRequestImage) UnmarshalBinary(b []byte) error {
	var res ImagesRequestImage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
