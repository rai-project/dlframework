package http

// This file contains authboss setup and necessary middlewares

import (
        "fmt"
        "context"
        "encoding/base64"
        // "encoding/json"
        "flag"
        // "fmt"
        // "html/template"
        // "io/ioutil"
        // "log"
        "net/http"
        // "os"
        // "regexp"
        // "strconv"
        // "time"

        // "github.com/BurntSushi/toml"
        "github.com/volatiletech/authboss"
        _ "github.com/volatiletech/authboss/auth"
        // "github.com/volatiletech/authboss/confirm"
        "github.com/volatiletech/authboss/defaults"
        // "github.com/volatiletech/authboss/lock"
        _ "github.com/volatiletech/authboss/logout"
        // aboauth "github.com/volatiletech/authboss/oauth2"
        // "github.com/volatiletech/authboss/otp/twofactor"
        // "github.com/volatiletech/authboss/otp/twofactor/sms2fa"
        // "github.com/volatiletech/authboss/otp/twofactor/totp2fa"
        // _ "github.com/volatiletech/authboss/recover"
        // _ "github.com/volatiletech/authboss/register"
        // "github.com/volatiletech/authboss/remember"

        "github.com/volatiletech/authboss-clientstate"
        // "github.com/volatiletech/authboss-renderer"

        // "golang.org/x/oauth2"
        // "golang.org/x/oauth2/google"

        "github.com/aarondl/tpl"
        // "github.com/go-chi/chi"
        "github.com/gorilla/schema"
        "github.com/gorilla/sessions"
        "github.com/justinas/nosurf"
)

var (
        flagDebug    = flag.Bool("debug", false, "output debugging information")
        flagDebugDB  = flag.Bool("debugdb", false, "output database on each request")
        flagDebugCTX = flag.Bool("debugctx", false, "output specific authboss related context keys on each request")
        flagAPI      = flag.Bool("api", true, "configure the app to be an api instead of an html app")
)

var (
        ab        = authboss.New()
        Ab        = ab
        database  = NewMemStorer()
        schemaDec = schema.NewDecoder()

        sessionStore abclientstate.SessionStorer
        cookieStore  abclientstate.CookieStorer

        templates tpl.Templates
)

const (
        sessionCookieName = "ab_mlmodelscope"
)

func setupAuthboss() {
        flag.Parse()

        // TODO: provide updated keys
        // Initialize Sessions and Cookies
        // Typically gorilla securecookie and sessions packages require
        // highly random secret keys that are not divulged to the public.
        //
        // In this example we use keys generated one time (if these keys ever become
        // compromised the gorilla libraries allow for key rotation, see gorilla docs)
        // The keys are 64-bytes as recommended for HMAC keys as per the gorilla docs.
        //
        // These values MUST be changed for any new project as these keys are already "compromised"
        // as they're in the public domain, if you do not change these your application will have a fairly
        // wide-opened security hole. You can generate your own with the code below, or using whatever method
        // you prefer:
        //
        //    func main() {
        //        fmt.Println(base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)))
        //    }
        //
        // We store them in base64 in the example to make it easy if we wanted to move them later to
        // a configuration environment var or file.
        cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
        sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
        cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
        sessionStore = abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
        cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false

        // ab.Config.Paths.RootURL = "https://www.mlmodelscope.org"
        ab.Config.Paths.RootURL = "http://localhost:8088"

        if !*flagAPI {
                // Prevent us from having to use Javascript in our basic HTML
                // to create a delete method, but don't override this default for the API
                // version
                ab.Config.Modules.LogoutMethod = "GET"
        }

        // Set up our server, session and cookie storage mechanisms.
        // These are all from this package since the burden is on the
        // implementer for these.
        ab.Config.Storage.Server = database
        ab.Config.Storage.SessionState = sessionStore
        ab.Config.Storage.CookieState = cookieStore


        // TODO: Set view and mail renderer
        ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}

        // The preserve fields are things we don't want to
        // lose when we're doing user registration (prevents having
        // to type them again)
        ab.Config.Modules.RegisterPreserveFields = []string{"email", "first_name", "last_name"}

        // TOTP2FAIssuer is the name of the issuer we use for totp 2fa
        ab.Config.Modules.TOTP2FAIssuer = "ABMlModelscope"
        ab.Config.Modules.RoutesRedirectOnUnauthed = true

        // This instantiates and uses every default implementation
        // in the Config.Core area that exist in the defaults package.
        // Just a convenient helper if you don't want to do anything fancy.
        defaults.SetCore(&ab.Config, *flagAPI, false)

        // Here we initialize the bodyreader as something customized in order to accept a name
        // parameter for our user as well as the standard e-mail and password.
        /* emailRule := defaults.Rules{ */
                // FieldName: "email", Required: true,
                // MatchError: "Must be a valid e-mail address",
                // MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]{1,}`),
        /* } */
        usernameRule := defaults.Rules{
                FieldName: "username", Required: true,
                MinLength: 2,
        }

        passwordRule := defaults.Rules{
                FieldName: "password", Required: true,
                MinLength: 4,
        }
        /* firstnameRule := defaults.Rules{ */
                // FieldName: "first_name", Required: true,
                // MinLength: 2,
        /* } */

        ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
                ReadJSON: *flagAPI,
                UseUsername: true,
                Rulesets: map[string][]defaults.Rules{
                        // "register":    {emailRule, passwordRule, firstnameRule},
                        "register":    {passwordRule, usernameRule},
                        "recover_end": {passwordRule},
                },
                // Confirms: map[string][]string{
                        // "register":    {"password", authboss.ConfirmPrefix + "password"},
                        // "recover_end": {"password", authboss.ConfirmPrefix + "password"},
                // },
                Whitelist: map[string][]string{
                        "register": []string{"email", "username", "password"},
                },
        }

        // Initialize authboss (instantiate modules etc.)SessionKey
        if err := ab.Init(); err != nil {
                panic(err)
        }

        // Setup router
        schemaDec.IgnoreUnknownKeys(true)
        // TODO: get access to our middleware
        // interface it with ab.LoadClientStateMiddleware
        // remember.Middleware(ab)
        // authboss.Middleware(..), lock.Middleware(ab), confirm.Middleware(ab)

        // TODO: add csrf token handling through routes
        // optionsHandler := func(w http.ResponseWriter, r *http.Request) {
        //         w.Header().Set("X-CSRF-TOKEN", nosurf.Token(r))
        //         w.WriteHeader(http.StatusOK)
        // }
        //
        // // We have to add each of the authboss get/post routes specifically because
        // // chi sees the 'Mount' above as overriding the '/*' pattern.
        // routes := []string{"login", "logout", "recover", "recover/end", "register"}
        // mux.MethodFunc("OPTIONS", "/*", optionsHandler)
        // for _, r := range routes {
        //         mux.MethodFunc("OPTIONS", "/auth/"+r, optionsHandler)
        // }

}

func nosurfing(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Failed to validate CSRF token:", nosurf.Reason(r))
		w.WriteHeader(http.StatusBadRequest)
	}))
	return surfing
}

func dataInjector(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := layoutData(w, &r)
		r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, data))
		handler.ServeHTTP(w, r)
	})
}

// layoutData is passing pointers to pointers be able to edit the current pointer
// to the request. This is still safe as it still creates a new request and doesn't
// modify the old one, it just modifies what we're pointing to in our methods so
// we're able to skip returning an *http.Request everywhere
func layoutData(w http.ResponseWriter, r **http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ab.LoadCurrentUser(r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).FirstName
	}

	return authboss.HTMLData{
		"loggedin":          userInter != nil,
		"current_user_name": currentUserName,
                // "csrf_token":        nosurf.Token(*r),
		"flash_success":     authboss.FlashSuccess(w, *r),
		"flash_error":       authboss.FlashError(w, *r),
	}
}


