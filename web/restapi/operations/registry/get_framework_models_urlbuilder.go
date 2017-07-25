package registry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"
)

// GetFrameworkModelsURL generates an URL for the get framework models operation
type GetFrameworkModelsURL struct {
	FrameworkName string

	FrameworkVersion *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetFrameworkModelsURL) WithBasePath(bp string) *GetFrameworkModelsURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetFrameworkModelsURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetFrameworkModelsURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/v1/framework/{framework_name}/models"

	frameworkName := o.FrameworkName
	if frameworkName != "" {
		_path = strings.Replace(_path, "{framework_name}", frameworkName, -1)
	} else {
		return nil, errors.New("FrameworkName is required on GetFrameworkModelsURL")
	}
	_basePath := o._basePath
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var frameworkVersion string
	if o.FrameworkVersion != nil {
		frameworkVersion = *o.FrameworkVersion
	}
	if frameworkVersion != "" {
		qs.Set("framework_version", frameworkVersion)
	}

	result.RawQuery = qs.Encode()

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetFrameworkModelsURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetFrameworkModelsURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetFrameworkModelsURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetFrameworkModelsURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetFrameworkModelsURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetFrameworkModelsURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
