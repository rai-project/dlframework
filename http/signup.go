package http

import (

  "github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
  "github.com/rai-project/dlframework/httpapi/restapi/operations/signup"
)

func SignupHandler(params signup.SignupParams) middleware.Responder {

	firstName := params.Body.FirstName
	lastName := params.Body.LastName
	username := params.Body.Username
	password := params.Body.Password
	affiliation := params.Body.Affiliation

	// check if first name/last name/affiliation matches in the DB
	// check if username already exists
	if (len(firstName) == 0 || len(lastName) == 0 || len(username) == 0 || len(password) == 0 || len(affiliation) == 0) {
    return NewError("Signup", errors.New("Incomplete Information"))
  }

	return signup.NewSignupOK()

}
