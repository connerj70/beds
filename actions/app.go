package actions

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	"beds/models"

	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.

func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_beds_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())

		app.Use(authenticate)

		var userResource UsersResource
		var bedsResource BedsResource
		app.Middleware.Skip(authenticate, HomeHandler, userResource.Create, userResource.New, userResource.SignIn, userResource.SignInPage, userResource.SignOut)

		app.GET("/", HomeHandler)
		app.GET("/index", Index)
		app.GET("/signin", userResource.SignInPage)
		app.POST("/signin", userResource.SignIn)
		app.GET("/signout", userResource.SignOut)
		app.POST("/users/by_email", userResource.FindByEmail)
		app.Resource("/users", userResource)
		app.POST("/beds/toggle_complete", bedsResource.ToggleComplete)
		app.Resource("/beds", bedsResource)
		app.POST("/friends/create", FriendsCreate)
		app.GET("/friends/list/{id}", FriendsList)
		app.GET("/friends/show", FriendsListPage)
		app.ServeFiles("/", assetsBox) // serve files from the public directory

		// Setup workers
		w := app.Worker
		w.Register("reset_daily_beds", func(worker.Args) error {

			if err := models.DB.RawQuery("UPDATE beds SET complete = false WHERE frequency = 1").Exec(); err != nil {
				return fmt.Errorf("failed to reset daily beds")
			}

			return nil
		})

		// Push jobs on worker at the end of every day
		go func() {
			t := time.NewTicker(1 * time.Minute)

			var file *os.File
			var err error
			var lastBedsResetTime time.Time
			_, err = os.Stat("/tmp/beds/reset_time")
			if err != nil {
				log.Println("reset_time file does not exist creating one now")
				file, err = os.Create("/tmp/beds/reset_time")
				if err != nil {
					log.Panicln("failed to create reset_time file: ", err)
				}
				now := time.Now()
				lastBedsResetTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 1, 0, time.UTC)
				currentTimeStr := lastBedsResetTime.Format(time.RFC3339)
				_, err := file.WriteAt([]byte(currentTimeStr), 0)
				if err != nil {
					log.Println("failed to write to reset_time file: ", err)
				}
			} else {
				file, err = os.Open("/tmp/beds/reset_time")
				if err != nil {
					log.Panicln("failed to open reset_time file: ", err)
				}

				fileContent, err := ioutil.ReadAll(file)
				if err != nil {
					log.Panicln("failed to read reset_time file: ", err)
				}

				lastBedsResetTime, err = time.Parse(time.RFC3339, string(fileContent))
				if err != nil {
					log.Panicln("failed to parse reset_time contents: ", err)
				}
			}

			for {
				select {
				case currentTime := <-t.C:
					currentTimeUTC := currentTime.UTC()
					if currentTimeUTC.After(lastBedsResetTime) {
						w.Perform(worker.Job{
							Queue:   "default",
							Handler: "reset_daily_beds",
						})

						nextResetTime := lastBedsResetTime.AddDate(0, 0, 1)
						nextResetTimeStr := nextResetTime.Format(time.RFC3339)
						_, err := file.WriteAt([]byte(nextResetTimeStr), 0)
						if err != nil {
							log.Println("failed to write to reset_time file: ", err)
						}
					}
				}
			}
		}()
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func authenticate(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		jwtString, err := c.Cookies().Get("jwt")
		if err != nil {
			return fmt.Errorf("missing authorization cookie")
		}
		if jwtString == "" {
			return c.Redirect(http.StatusUnauthorized, "/signin")
		}
		c.Set("token", jwtString)
		// Parse the jwt
		_, err = jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(envy.Get("BEDS_JWT_SECRET", "")), nil
		})
		if err != nil {
			return fmt.Errorf("there was in issue with the jwt: %w", err)
		}
		// If the user is signed in, call the next handler
		return next(c)
	}
}
