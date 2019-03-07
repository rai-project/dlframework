// Code generated by go-swagger; DO NOT EDIT.

package predict

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new predict API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for predict API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
Close closes a predictor clear it s memory
*/
func (a *Client) Close(params *CloseParams) (*CloseOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCloseParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "Close",
		Method:             "POST",
		PathPattern:        "/predict/close",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &CloseReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*CloseOK), nil

}

/*
Dataset datasets method receives a single dataset and runs the predictor on all elements of the dataset the result is a prediction feature list
*/
func (a *Client) Dataset(params *DatasetParams) (*DatasetOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDatasetParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "Dataset",
		Method:             "POST",
		PathPattern:        "/predict/dataset",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &DatasetReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DatasetOK), nil

}

/*
DatasetStream datasets method receives a single dataset and runs the predictor on all elements of the dataset the result is a prediction feature stream
*/
func (a *Client) DatasetStream(params *DatasetStreamParams) (*DatasetStreamOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDatasetStreamParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DatasetStream",
		Method:             "POST",
		PathPattern:        "/predict/stream/dataset",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &DatasetStreamReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DatasetStreamOK), nil

}

/*
Images images method receives a list of base64 encoded images and runs the predictor on all the images the result is a prediction feature list for each image
*/
func (a *Client) Images(params *ImagesParams) (*ImagesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewImagesParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "Images",
		Method:             "POST",
		PathPattern:        "/predict/images",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ImagesReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ImagesOK), nil

}

/*
ImagesStream images stream method receives a list base64 encoded images and runs the predictor on all the images the result is a prediction feature stream for each image
*/
func (a *Client) ImagesStream(params *ImagesStreamParams) (*ImagesStreamOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewImagesStreamParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "ImagesStream",
		Method:             "POST",
		PathPattern:        "/predict/stream/images",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ImagesStreamReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ImagesStreamOK), nil

}

/*
Open opens a predictor and returns an id where the predictor is accessible the id can be used to perform inference requests
*/
func (a *Client) Open(params *OpenParams) (*OpenOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewOpenParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "Open",
		Method:             "POST",
		PathPattern:        "/predict/open",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &OpenReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*OpenOK), nil

}

/*
Reset resets method clears the internal cache of the predictors
*/
func (a *Client) Reset(params *ResetParams) (*ResetOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewResetParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "Reset",
		Method:             "POST",
		PathPattern:        "/predict/reset",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ResetReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ResetOK), nil

}

/*
Urls urls method receives a list of urls and runs the predictor on all the urls the result is a list of predicted features for all the urls
*/
func (a *Client) Urls(params *UrlsParams) (*UrlsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUrlsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "URLs",
		Method:             "POST",
		PathPattern:        "/predict/urls",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &UrlsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UrlsOK), nil

}

/*
UrlsStream urls stream method receives a stream of urls and runs the predictor on all the urls the result is a prediction feature stream for each url
*/
func (a *Client) UrlsStream(params *UrlsStreamParams) (*UrlsStreamOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUrlsStreamParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "URLsStream",
		Method:             "POST",
		PathPattern:        "/predict/stream/urls",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &UrlsStreamReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UrlsStreamOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
