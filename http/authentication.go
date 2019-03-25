package http

import (
        // "fmt"
        // "net/http"
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/authentication"
	"github.com/rai-project/passlib"
        // "github.com/volatiletech/authboss"
)

// dummy for MongoDB User Database
type userRecord struct {
	firstname   string
	lastname    string
	username    string
	password    string
	affiliation string
}

var userTable = map[string]*userRecord{}

func CheckPassword(username string, password string) bool {

	hash := userTable[username].password

	newHash, err := passlib.Verify(password, hash)
	if err != nil {
		// incorrect password, malformed hash, etc.
		// either way, reject
		return false
	}

	// The context has decided, as per its policy, that
	// the hash which was used to validate the password
	// should be changed. It has upgraded the hash using
	// the verified password.
	if newHash != "" {
		userTable[username].password = newHash
	}

	return true

}

func LoginHandler(params authentication.LoginParams) middleware.Responder {

	/* username := params.Body.Username */
	// password := params.Body.Password
        //
	// if val, ok := userTable[username]; ok {
	//         if CheckPassword(val.username, password) == false {
	//                 return NewError("Login", errors.New("Incorrect credentials"))
	//         }
	// } else {
	//         return NewError("Login", errors.New("Not signed up!"))
	/* } */
        /* ctx := params.HTTPRequest.Context() */
        // ab_instance := ctx.Value("authboss_instance").(*authboss.Authboss)
        // abuser := ab_instance.CurrentUserP(params.HTTPRequest)
        /* fmt.Printf(abuser.GetPID()) */
        /* params.HTTPRequest */
        // fmt.Printf("LoginHandler")
        // a := &Auth{ab}
        /* a.LoginGet() */

	return authentication.NewLoginOK()

}

func GenerateHash(password string) string {

	hash, err := passlib.Hash(password)
	if err != nil {
		// couldn't hash password for some reason
		return "xxx"
	}

	return hash

}

func SignupHandler(params authentication.SignupParams) middleware.Responder {

	firstName := params.Body.FirstName
	lastName := params.Body.LastName
	username := params.Body.Username
	password := params.Body.Password
	affiliation := params.Body.Affiliation

	if len(firstName) == 0 || len(lastName) == 0 || len(username) == 0 || len(password) == 0 || len(affiliation) == 0 {
		return NewError("Signup", errors.New("Incomplete Information"))
	}

	encryption := GenerateHash(password)
	if encryption == "xxx" {
		return NewError("Signup", errors.New("Could not generate password"))
	}

	userTable[username] = &userRecord{firstname: firstName, lastname: lastName, username: username, password: encryption, affiliation: affiliation}

	return authentication.NewSignupOK()

}
