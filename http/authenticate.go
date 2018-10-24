package http

import (
        "context"
        "encoding/base64"
        "encoding/json"
        "flag"
        "fmt"
        "html/template"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "regexp"
        "strconv"
        "time"

        "github.com/BurntSushi/toml"
        "github.com/volatiletech/authboss"
        _ "github.com/volatiletech/authboss/auth"
        "github.com/volatiletech/authboss/confirm"
        "github.com/volatiletech/authboss/defaults"
        "github.com/volatiletech/authboss/lock"
        _ "github.com/volatiletech/authboss/logout"
        aboauth "github.com/volatiletech/authboss/oauth2"
        "github.com/volatiletech/authboss/otp/twofactor"
        "github.com/volatiletech/authboss/otp/twofactor/sms2fa"
        "github.com/volatiletech/authboss/otp/twofactor/totp2fa"
        _ "github.com/volatiletech/authboss/recover"
        _ "github.com/volatiletech/authboss/register"
        "github.com/volatiletech/authboss/remember"

        "github.com/volatiletech/authboss-clientstate"
        "github.com/volatiletech/authboss-renderer"

        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"

        "github.com/aarondl/tpl"
        "github.com/go-chi/chi"
        "github.com/gorilla/schema"
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

	ab.Config.Paths.RootURL = "https://www.mlmodelscope.org"

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
        emailRule := defaults.Rules{
                FieldName: "email", Required: true,
                MatchError: "Must be a valid e-mail address",
                MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]{1,}`),
        }
        passwordRule := defaults.Rules{
                FieldName: "password", Required: true,
                MinLength: 4,
        }
        firstnameRule := defaults.Rules{
                FieldName: "first_name", Required: true,
                MinLength: 2,
        }

        ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
                ReadJSON: *flagAPI,
                Rulesets: map[string][]defaults.Rules{
                        "register":    {emailRule, passwordRule, firstnameRule},
                        "recover_end": {passwordRule},
                },
                Confirms: map[string][]string{
                        "register":    {"password", authboss.ConfirmPrefix + "password"},
                        "recover_end": {"password", authboss.ConfirmPrefix + "password"},
                },
                Whitelist: map[string][]string{
                        "register": []string{"email", "first_name", "password"},
                },
        }

        oauthcreds := struct {
                ClientID     string `toml:"client_id"`
                ClientSecret string `toml:"client_secret"`
        }{}

	// Set up 2fa
        twofaRecovery := &twofactor.Recovery{Authboss: ab}
        if err := twofaRecovery.Setup(); err != nil {
                panic(err)
        }

        totp := &totp2fa.TOTP{Authboss: ab}
        if err := totp.Setup(); err != nil {
                panic(err)
        }

        sms := &sms2fa.SMS{Authboss: ab, Sender: smsLogSender{}}
        if err := sms.Setup(); err != nil {
                panic(err)
        }

	// Set up Google OAuth2 if we have credentials in the
        // file oauth2.toml for it.
        _, err := toml.DecodeFile("oauth2.toml", &oauthcreds)
        if err == nil && len(oauthcreds.ClientID) != 0 && len(oauthcreds.ClientSecret) != 0 {
                fmt.Println("oauth2.toml exists, configuring google oauth2")
                ab.Config.Modules.OAuth2Providers = map[string]authboss.OAuth2Provider{
                        "google": authboss.OAuth2Provider{
                                OAuth2Config: &oauth2.Config{
                                        ClientID:     oauthcreds.ClientID,
                                        ClientSecret: oauthcreds.ClientSecret,
                                        Scopes:       []string{`profile`, `email`},
                                        Endpoint:     google.Endpoint,
                                },
                                FindUserDetails: aboauth.GoogleUserDetails,
                        },
                }
        } else if os.IsNotExist(err) {
                fmt.Println("oauth2.toml doesn't exist, not registering oauth2 handling")
        } else {
                fmt.Println("error loading oauth2.toml:", err)
        }

	// Initialize authboss (instantiate modules etc.)
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

}


