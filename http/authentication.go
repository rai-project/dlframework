package http

import (
        "io/ioutil"
        "bytes"
        "encoding/json"
        "net/http"
        "github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/authentication"
        auth "github.com/volatiletech/authboss/auth"
        register "github.com/volatiletech/authboss/register"
        logout "github.com/volatiletech/authboss/logout"
)

func LoginHandler(params authentication.LoginParams, principal *models.User) middleware.Responder {
        return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer){
                a := &auth.Auth{ab}
                req := params.HTTPRequest
                requestByte, _ := json.Marshal(principal)
                req.Body = ioutil.NopCloser(bytes.NewReader(requestByte))

                a.LoginPost(rw, req)
        })

}

func SignupHandler(params authentication.SignupParams) middleware.Responder {
        return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer){
                r := &register.Register{ab}
                req := params.HTTPRequest
                requestByte, _ := json.Marshal(params.Body)
                req.Body = ioutil.NopCloser(bytes.NewReader(requestByte))
                r.Post(rw, req)
        })
}

func UserInfoHandler(params authentication.UserInfoParams) middleware.Responder {
        userInter, err := ab.LoadCurrentUser(&params.HTTPRequest)
        if userInter == nil || err != nil {
            return authentication.NewUserInfoOK().
            WithPayload(&models.DlframeworkUserInfoResponse{
                    Outcome: "fail",
            })
        } else {
            return authentication.NewUserInfoOK().
            WithPayload(&models.DlframeworkUserInfoResponse{
                    Outcome: "success",
                    Username: userInter.(*User).Username,
                    Email: userInter.(*User).Email,
                    FirstName: userInter.(*User).FirstName,
                    LastName: userInter.(*User).LastName,
            })
        }
}

func LogoutHandler(params authentication.LogoutParams) middleware.Responder {
        return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer){
                l := &logout.Logout{ab}
                req := params.HTTPRequest
                l.Logout(rw, req)
        })
}
