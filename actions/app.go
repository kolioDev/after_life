package actions

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/kolioDev/after_life/models"
	"github.com/markbates/goth/gothic"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

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
			Env: ENV,
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_after_life_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		app.GET("/", HomeHandler)

		oauth := app.Group("/oauth")
		oauth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
		oauth.GET("/{provider}/callback", AuthOAuthCallback)

		auth := app.Group("/auth")
		auth.GET("/access/{session_identifier}", AuthGetInitialTokens)
		auth.POST("/token", AuthRefreshToken)
		auth.POST("/reset", AuthResetSignUp)

		user := app.Group("/user")
		user.GET("/keys", AuthMiddleware(UsersGetConfirmationKeys, true))
		user.POST("/confirm/keys", AuthMiddleware(UsersConfirmKeysSeen, true))

		app.ANY("/graphql", AuthMiddleware(GraphqlIndex, false))

	}

	return app
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

//Fixes bug - Unable to read post values
//Copy c.Request().Body into c.Request().PostForm
func postParse(c buffalo.Context) (buffalo.Context, error) {
	//Only works wor POST/PUT request method
	if c.Request().Method != "POST" && c.Request().Method != "PUT" {
		return c, nil
	}

	var vs map[string]interface{}
	err := json.Unmarshal([]byte(fmt.Sprint(c.Request().Body)), &vs)
	if err != nil {
		return c, errors.WithStack(err)
	}

	vals := url.Values{}
	for k, v := range vs {
		vals.Add(k, fmt.Sprint(v))
	}
	c.Request().PostForm = vals
	return c, nil
}
