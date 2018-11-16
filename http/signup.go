package http

import (

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/signup"
	"github.com/k0kubun/pp"
	//"github.com/rai-project/passlib"
	//"github.com/volatiletech/authboss"
        //_ "github.com/volatiletech/authboss/auth"
	//register "github.com/volatiletech/authboss/register"
)


// dummy for MongoDB User Database
/*type userRecord struct {
	firstname		string
	lastname		string
	username		string
	password		string
	affiliation		string
	email			string
}
var userTable = map[string]*userRecord{}

func GenerateHash(password string) string {

	hash, err := passlib.Hash(password)
	if err != nil {
		// couldn't hash password for some reason
		return "xxx"
	}

	return hash

}*/

func SignupHandler(params signup.SignupParams) middleware.Responder {

	firstName := params.Body.FirstName
	lastName := params.Body.LastName
	username := params.Body.Username
	password := params.Body.Password
	affiliation := params.Body.Affiliation
	email := params.Body.Email

	if (len(firstName) == 0 || len(lastName) == 0 || len(username) == 0 || len(password) == 0 || len(affiliation) == 0 || len(email) == 0) {
		return NewError("Signup", errors.New("Incomplete Information"))
	}

	pp.Println("ABCDCDCDCDCDCCDCD")
	//encryption := GenerateHash(password)
	//if encryption == "xxx" {
	//	return NewError("Signup", errors.New("Could not generate password"))
	//}

	//ctx := params.HTTPRequest.Context()
	//userTable[username] = &userRecord{firstname: firstName, lastname: lastName, username: username, password: encryption, affiliation: affiliation, email: email}

	// interface with storer
	// fetch ab instance
	//ab_instance := ctx.Value("authboss_instance").(*authboss.Authboss)

	// fetch user instance/current user from ab
	//abuser := ab_instance.CurrentUserP(params.HTTPRequest)
	// refer to example to see what else we can do for authentication
	//abuser := ab.CurrentUserP(r)
	//user := abuser.(*User)

	// Register user
	//reg := &register.Register{ab_instance}
	//reg.Post(params.HTTPResponseWriter, params.HTTPRequest)

	// initiate register library at the mlmodelscope/pkg end
	// and probably pass it through context/or 
	return signup.NewSignupOK()

}
