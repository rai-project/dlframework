// Code generated by go-swagger; DO NOT EDIT.

package predict

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/rai-project/dlframework/httpapi/models"
)

// DatasetReader is a Reader for the Dataset structure.
type DatasetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DatasetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewDatasetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDatasetOK creates a DatasetOK with default headers values
func NewDatasetOK() *DatasetOK {
	return &DatasetOK{}
}

/*DatasetOK handles this case with default header values.

A successful response.
*/
type DatasetOK struct {
	Payload *models.DlframeworkFeaturesResponse
}

func (o *DatasetOK) Error() string {
	return fmt.Sprintf("[POST /predict/dataset][%d] datasetOK  %+v", 200, o.Payload)
}

func (o *DatasetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DlframeworkFeaturesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
