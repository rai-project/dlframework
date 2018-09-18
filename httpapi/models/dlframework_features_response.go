// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// DlframeworkFeaturesResponse dlframework features response
// swagger:model dlframeworkFeaturesResponse
type DlframeworkFeaturesResponse struct {

	// id
	ID string `json:"id,omitempty"`

	// responses
	Responses []*DlframeworkFeatureResponse `json:"responses"`

	// trace id
	TraceID *DlframeworkTraceID `json:"trace_id,omitempty"`
}

// Validate validates this dlframework features response
func (m *DlframeworkFeaturesResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateResponses(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTraceID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DlframeworkFeaturesResponse) validateResponses(formats strfmt.Registry) error {

	if swag.IsZero(m.Responses) { // not required
		return nil
	}

	for i := 0; i < len(m.Responses); i++ {
		if swag.IsZero(m.Responses[i]) { // not required
			continue
		}

		if m.Responses[i] != nil {
			if err := m.Responses[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("responses" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *DlframeworkFeaturesResponse) validateTraceID(formats strfmt.Registry) error {

	if swag.IsZero(m.TraceID) { // not required
		return nil
	}

	if m.TraceID != nil {
		if err := m.TraceID.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("trace_id")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DlframeworkFeaturesResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DlframeworkFeaturesResponse) UnmarshalBinary(b []byte) error {
	var res DlframeworkFeaturesResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
