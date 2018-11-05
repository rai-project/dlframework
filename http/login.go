package http

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	"github.com/rai-project/passlib"
)

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

func LoginHandler(params login.LoginParams) middleware.Responder {

	username := params.Body.Username
	password := params.Body.Password

	if val, ok := userTable[username]; ok {
		if CheckPassword(val.username, password) == false {
			return NewError("Login", errors.New("Incorrect credentials"))
		}
	} else {
		return NewError("Login", errors.New("Not signed up!"))
	}

	return login.NewLoginOK()

}
