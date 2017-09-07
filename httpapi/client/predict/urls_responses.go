// Code generated by go-swagger; DO NOT EDIT.

package predict

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/rai-project/dlframework/httpapi/models"
)

// UrlsReader is a Reader for the Urls structure.
type UrlsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UrlsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewUrlsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUrlsOK creates a UrlsOK with default headers values
func NewUrlsOK() *UrlsOK {
	return &UrlsOK{}
}

/*UrlsOK handles this case with default header values.

UrlsOK urls o k
*/
type UrlsOK struct {
	Payload *models.DlframeworkFeaturesResponse
}

func (o *UrlsOK) Error() string {
	return fmt.Sprintf("[POST /v1/predict/urls][%d] urlsOK  %+v", 200, o.Payload)
}

func (o *UrlsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DlframeworkFeaturesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
