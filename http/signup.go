package http

import (

  "github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
  "github.com/rai-project/dlframework/httpapi/restapi/operations/signup"
	"github.com/rai-project/passlib"
)

// dummy for MongoDB User Database
type userRecord struct {
	firstname		string
	lastname		string
	username		string
	password		string
	affiliation	string
}
var userTable = map[string]*userRecord{}

func GenerateHash(password string) string {

	hash, err := passlib.Hash(password)
	if err != nil {
		// couldn't hash password for some reason
		return "xxx"
	}

	return hash

}

func SignupHandler(params signup.SignupParams) middleware.Responder {

	firstName := params.Body.FirstName
	lastName := params.Body.LastName
	username := params.Body.Username
	password := params.Body.Password
	affiliation := params.Body.Affiliation

	if (len(firstName) == 0 || len(lastName) == 0 || len(username) == 0 || len(password) == 0 || len(affiliation) == 0) {
    return NewError("Signup", errors.New("Incomplete Information"))
  }

	encryption := GenerateHash(password)
	if encryption == "xxx" {
		return NewError("Signup", errors.New("Could not generate password"))
	}

	userTable[username] = &userRecord{firstname: firstName, lastname: lastName, username: username, password: encryption, affiliation: affiliation}

	return signup.NewSignupOK()

}
