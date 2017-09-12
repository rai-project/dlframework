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

// DatasetStreamReader is a Reader for the DatasetStream structure.
type DatasetStreamReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DatasetStreamReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewDatasetStreamOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDatasetStreamOK creates a DatasetStreamOK with default headers values
func NewDatasetStreamOK() *DatasetStreamOK {
	return &DatasetStreamOK{}
}

/*DatasetStreamOK handles this case with default header values.

(streaming responses)
*/
type DatasetStreamOK struct {
	Payload *models.DlframeworkFeatureResponse
}

func (o *DatasetStreamOK) Error() string {
	return fmt.Sprintf("[POST /v1/predict/stream/dataset][%d] datasetStreamOK  %+v", 200, o.Payload)
}

func (o *DatasetStreamOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DlframeworkFeatureResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}