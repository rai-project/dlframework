package http

import (

  "github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
  "github.com/rai-project/dlframework/httpapi/restapi/operations/login"
)

func LoginHandler(params login.LoginParams) middleware.Responder {

	username := params.Body.Username
	password := params.Body.Password

	// verify the information by querying the database
	if (username != "admin" || password != "admin") {
		return NewError("Login", errors.New("Incorrect credentials"))
	}

	return login.NewLoginOK()

}
