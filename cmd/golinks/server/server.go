package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/popodidi/log"

	"github.com/haostudio/golinks/cmd/golinks/server/modules/authapi"
	"github.com/haostudio/golinks/cmd/golinks/server/modules/authweb"
	"github.com/haostudio/golinks/cmd/golinks/server/modules/landingweb"
	"github.com/haostudio/golinks/cmd/golinks/server/modules/linkapi"
	"github.com/haostudio/golinks/cmd/golinks/server/modules/linkweb"
	"github.com/haostudio/golinks/cmd/golinks/server/modules/redirect"
	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/link"
)

// Config defines the golink server config.
type Config struct {
	Gin          *gin.Engine
	Traced       bool
	AuthProvider auth.Provider
	LinkStore    link.Store
}

// New returns a new server instance.
func New(conf Config) http.Handler {
	// Server logger
	logger := log.New("server")

	// Create our HTTP Router.
	router := conf.Gin

	// Configure HTTP Router Settings
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = false
	router.HandleMethodNotAllowed = false
	router.ForwardedByClientIP = true
	router.AppEngine = false
	router.UseRawPath = false
	router.UnescapePathValues = true

	// Log server config
	logger.Info("server link store: %s", conf.LinkStore)
	logger.Info("server auth provider: %s", conf.AuthProvider)

	// Setup middlewares.
	if conf.Traced {
		router.Use(middlewares.Trace)
	}
	router.Use(middlewares.CtxLogger)
	router.Use(middlewares.PanicCatcher)

	// nolint: godox
	// FIXME: add favicon and remove this hack
	router.Use(func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.URL.String(), "/favicon.ico") {
			ctx.AbortWithStatus(http.StatusNotFound)
		}
	})

	// Landing page
	landingweb.Register(router, landingweb.Config{
		Traced: conf.Traced,
	})

	// Link module
	lnGroup := router.Group("links")
	lnGroup.Use(middlewares.Auth(conf.AuthProvider))
	linkweb.Register(lnGroup, linkweb.Config{
		Store:  conf.LinkStore,
		Traced: conf.Traced,
	})

	// Org web module
	orgGroup := router.Group("org")
	authweb.Register(orgGroup, authweb.Config{
		Traced:   conf.Traced,
		Provider: conf.AuthProvider,
	})

	// Link api module
	lnAPIGroup := router.Group("api/links")
	lnAPIGroup.Use(middlewares.Auth(conf.AuthProvider))
	linkapi.Register(lnAPIGroup, conf.LinkStore)

	// Auth module
	authGroup := router.Group("api/orgs")
	authGroup.Use(middlewares.Auth(conf.AuthProvider))
	authapi.Register(authGroup, conf.AuthProvider)

	// Use redirect handler by default.
	router.NoRoute(
		middlewares.Auth(conf.AuthProvider),
		middlewares.OrgRateLimit(5, time.Second),
		redirect.Handler(redirect.Config{
			Traced: conf.Traced,
			Store:  conf.LinkStore,
		}),
	)

	return router
}
