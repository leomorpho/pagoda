package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/middleware"
	"github.com/mikestefanello/pagoda/pkg/services"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

const (
	routeNameForgotPassword       = "forgot_password"
	routeNameForgotPasswordSubmit = "forgot_password.submit"
	routeNameLogin                = "login"
	routeNameLoginSubmit          = "login.submit"
	routeNameLogout               = "logout"
	routeNameRegister             = "register"
	routeNameRegisterSubmit       = "register.submit"
	routeNameResetPassword        = "reset_password"
	routeNameResetPasswordSubmit  = "reset_password.submit"
	routeNameVerifyEmail          = "verify_email"
	routeNameContact              = "contact"
	routeNameContactSubmit        = "contact.submit"
	routeNameAbout                = "about"
	routeNameHome                 = "home"
	routeNameSearch               = "search"
	routeNameDashboard            = "dashboard"
)

// BuildRouter builds the router
func BuildRouter(c *services.Container) {
	// Create a slog logger, which:
	//   - Logs to json.
	// TODO: add option to log to file
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Static files with proper cache control
	// funcmap.File() should be used in templates to append a cache key to the URL in order to break cache
	// after each server restart
	c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)).
		Static(config.StaticPrefix, config.StaticDir)

	// Custom handler for serving the service worker script with specific headers
	c.Web.GET("/service-worker.js", func(ctx echo.Context) error {
		// Set headers to allow the service worker scope to be at the root level
		ctx.Response().Header().Set("Service-Worker-Allowed", "/")
		// Set caching headers - adjust max-age as needed
		ctx.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", c.Config.Cache.Expiration.StaticFile))

		return ctx.File("./service-worker.js")
	})

	// Non static file route group
	g := c.Web.Group("")

	// Force HTTPS, if enabled
	if c.Config.HTTP.TLS.Enabled {
		g.Use(echomw.HTTPSRedirect())
	}

	g.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover(),
		echomw.Secure(),
		echomw.RequestID(),
		echomw.Gzip(),
		slogecho.New(logger),
		middleware.LogRequestID(),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Timeout: c.Config.App.Timeout,
		}),
		session.Middleware(sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))),
		middleware.LoadAuthenticatedUser(c.Auth),
		middleware.ServeCachedPage(c.Cache),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup: "form:csrf,header:X-CSRF-Token,query:csrf",
		}),
	)

	// Base controller
	ctr := controller.NewController(c)

	// Error handler
	err := errorHandler{Controller: ctr}
	c.Web.HTTPErrorHandler = err.Get

	// Example routes
	navRoutes(c, g, ctr)
	userRoutes(c, g, ctr)
}

func navRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {
	home := home{Controller: ctr}
	g.GET("/", home.Get).Name = routeNameHome

	search := search{Controller: ctr}
	g.GET("/search", search.Get).Name = routeNameSearch

	about := about{Controller: ctr}
	g.GET("/about", about.Get).Name = routeNameAbout

	contact := contact{Controller: ctr}
	g.GET("/contact", contact.Get).Name = routeNameContact
	g.POST("/contact", contact.Post).Name = routeNameContactSubmit
}

func userRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {
	auth := g.Group("/auth", middleware.RequireAuthentication())

	logout := logout{Controller: ctr}
	auth.GET("/logout", logout.Get, middleware.RequireAuthentication()).Name = routeNameLogout

	dashboard := dashboard{Controller: ctr}
	auth.GET("/dashboard", dashboard.Get, middleware.RequireAuthentication()).Name = routeNameDashboard

	verifyEmail := verifyEmail{Controller: ctr}
	g.GET("/email/verify/:token", verifyEmail.Get).Name = routeNameVerifyEmail

	noAuth := g.Group("/user", middleware.RequireNoAuthentication())
	login := login{Controller: ctr}
	noAuth.GET("/login", login.Get).Name = routeNameLogin
	noAuth.POST("/login", login.Post).Name = routeNameLoginSubmit

	register := register{Controller: ctr}
	noAuth.GET("/register", register.Get).Name = routeNameRegister
	noAuth.POST("/register", register.Post).Name = routeNameRegisterSubmit

	forgot := forgotPassword{Controller: ctr}
	noAuth.GET("/password", forgot.Get).Name = routeNameForgotPassword
	noAuth.POST("/password", forgot.Post).Name = routeNameForgotPasswordSubmit

	resetGroup := noAuth.Group("/password/reset",
		middleware.LoadUser(c.ORM),
		middleware.LoadValidPasswordToken(c.Auth),
	)
	reset := resetPassword{Controller: ctr}
	resetGroup.GET("/token/:user/:password_token/:token", reset.Get).Name = routeNameResetPassword
	resetGroup.POST("/token/:user/:password_token/:token", reset.Post).Name = routeNameResetPasswordSubmit

}
