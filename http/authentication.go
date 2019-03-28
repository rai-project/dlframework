package http

import (
        "fmt"
        "io/ioutil"
        "bytes"
        "encoding/json"
        "net/http"
        "github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	// "github.com/pkg/errors"
	"github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/authentication"
        "github.com/volatiletech/authboss"
        auth "github.com/volatiletech/authboss/auth"
        register "github.com/volatiletech/authboss/register"
        logout "github.com/volatiletech/authboss/logout"
        "github.com/k0kubun/pp"
)

func LoginHandler(params authentication.LoginParams, principal *models.User) middleware.Responder {
        return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer){
                a := &auth.Auth{ab}
                req := params.HTTPRequest
                pp.Println(authboss.GetSession(req, authboss.SessionKey))
                pp.Println(principal)
                // b, err := ioutil.ReadAll(req.Body)
                requestByte, _ := json.Marshal(principal)
                fmt.Println(string(requestByte))
                req.Body = ioutil.NopCloser(bytes.NewReader(requestByte))

                pp.Println("Login")
                a.LoginPost(rw, req)
                // pp.Println(req.Context().Value(authboss.CTXKeyUser))
                // u = ab.CurrentUser(req)
                // fmt.Println(u.GetPID())
        })
	// return authentication.NewLoginOK()

}

func SignupHandler(params authentication.SignupParams) middleware.Responder {
        return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer){
                r := &register.Register{ab}
                req := params.HTTPRequest
                requestByte, _ := json.Marshal(params.Body)
                fmt.Println(string(requestByte))
                req.Body = ioutil.NopCloser(bytes.NewReader(requestByte))
                r.Post(rw, req)
        })
	// return authentication.NewSignupOK()
}

func UserInfoHandler(params authentication.UserInfoParams) middleware.Responder {
        u, _ := ab.CurrentUser(params.HTTPRequest)
        pp.Println(u)
        return authentication.NewUserInfoOK().
        WithPayload(&models.DlframeworkUserInfoResponse{
                // Email: u.GetEmail(),
                Username: u.GetPID(),
        })
}

func LogoutHandler(params authentication.LogoutParams) middleware.Responder {
        return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer){
                l := &logout.Logout{ab}
                req := params.HTTPRequest
                l.Logout(rw, req)
        })
}
