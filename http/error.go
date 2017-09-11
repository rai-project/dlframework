package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"

	rerrors "github.com/go-openapi/errors"
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
	name := fmt.Sprintf("\"name\": \"%v\"", e.name)
	message := fmt.Sprintf("\"message\": \"%v\"", e.message)
	code := fmt.Sprintf("\"code\": %v", e.code)

	var stack string
	stackData := strings.Split(fmt.Sprintf("%+v", errors.WithStack(e.message)), "\n")
	bts, err := json.Marshal(stackData)
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
	default:
		rerrors.ServeError(rw, r, err)
	}
}
