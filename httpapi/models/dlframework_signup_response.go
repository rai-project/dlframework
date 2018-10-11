// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// DlframeworkSignupResponse dlframework signup response
// swagger:model dlframeworkSignupResponse
type DlframeworkSignupResponse struct {

	// outcome
	Outcome string `json:"outcome,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this dlframework signup response
func (m *DlframeworkSignupResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DlframeworkSignupResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DlframeworkSignupResponse) UnmarshalBinary(b []byte) error {
	var res DlframeworkSignupResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
