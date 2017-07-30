package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	runtimeerrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"
)

type Error struct {
	message  error
	code     int
	name     string
	rw       http.ResponseWriter
	producer runtime.Producer
}

func NewError(name string, message error) *Error {
	return &Error{
		message: message,
		code:    http.StatusBadRequest,
		name:    name,
	}
}

func (e *Error) MarshalJSON() ([]byte, error) {
	name := fmt.Sprintf("\"name\": \"%s\"", e.name)
	message := fmt.Sprintf("\"message\": \"%s\"", e.message.Error())
	code := fmt.Sprintf("\"code\": %d", e.code)
	stackData := strings.Split(fmt.Sprintf("%+v", errors.WithStack(e.message)), "\n")
	bts, err := json.Marshal(stackData)

	var stack string
	if err != nil {
		stack = fmt.Sprintf("\"stack\": []")
	} else {
		stack = fmt.Sprintf("\"stack\": %s", string(bts))
	}
	res := fmt.Sprintf("{%s, %s, %s, %s}", name, message, code, stack)
	return []byte(res), nil
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s:%v", e.name, e.message.Error())
}

func (e *Error) WithCode(code int) *Error {
	e.code = code
	return e
}

func (e *Error) WithName(name string) *Error {
	e.name = name
	return e
}

func (e *Error) WithMessage(message error) *Error {
	e.message = message
	return e
}

func (e *Error) Code() int32 {
	return int32(e.code)
}

func (e *Error) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(e.code)
	producer.Produce(rw, e)
}

// ServeError the error handler interface implemenation
func ServeError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.Header().Set("Content-Type", "application/json")
	switch e := err.(type) {
	case *Error:
		e.WriteResponse(rw, runtime.JSONProducer())
	case *runtimeerrors.CompositeError:
		er := flattenComposite(e)
		ServeError(rw, r, er.Errors[0])
	case *runtimeerrors.MethodNotAllowedError:
		rw.Header().Add("Allow", strings.Join(err.(*runtimeerrors.MethodNotAllowedError).Allowed, ","))
		rw.WriteHeader(asHTTPCode(int(e.Code())))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	case runtimeerrors.Error:
		if e == nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write(errorAsJSON(runtimeerrors.New(http.StatusInternalServerError, "Unknown error")))
			return
		}
		rw.WriteHeader(asHTTPCode(int(e.Code())))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	default:
		rw.WriteHeader(http.StatusInternalServerError)
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(runtimeerrors.New(http.StatusInternalServerError, err.Error())))
		}
	}
}

func flattenComposite(errs *runtimeerrors.CompositeError) *runtimeerrors.CompositeError {
	var res []error
	for _, er := range errs.Errors {
		switch e := er.(type) {
		case *runtimeerrors.CompositeError:
			if len(e.Errors) > 0 {
				flat := flattenComposite(e)
				if len(flat.Errors) > 0 {
					res = append(res, flat.Errors...)
				}
			}
		default:
			if e != nil {
				res = append(res, e)
			}
		}
	}
	return runtimeerrors.CompositeValidationError(res...)
}

func errorAsJSON(err runtimeerrors.Error) []byte {
	b, _ := json.Marshal(struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}{err.Code(), err.Error()})
	return b
}

func asHTTPCode(input int) int {
	if input >= 600 {
		return 422
	}
	return input
}
