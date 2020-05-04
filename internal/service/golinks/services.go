package golinks

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/popodidi/log"
	"github.com/soheilhy/cmux"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/link"
	"github.com/haostudio/golinks/internal/service"
	"github.com/haostudio/golinks/internal/service/golinks/modules/authapi"
	"github.com/haostudio/golinks/internal/service/golinks/modules/authweb"
	"github.com/haostudio/golinks/internal/service/golinks/modules/landingweb"
	"github.com/haostudio/golinks/internal/service/golinks/modules/linkapi"
	"github.com/haostudio/golinks/internal/service/golinks/modules/linkweb"
	"github.com/haostudio/golinks/internal/service/golinks/modules/redirect"
)

// Config defines the golinks http service config.
type Config struct {
	Gin          *gin.Engine
	Address      string
	Traced       bool
	AuthProvider auth.Provider
	LinkStore    link.Store
}

// New returns a golinks http service.
func New(config Config) service.Service {
	return &svc{
		Config: config,
	}
}

type svc struct {
	Config
}

func (s *svc) String() string {
	return "golinks.http"
}

func (s *svc) Matchers() []cmux.Matcher {
	return []cmux.Matcher{cmux.HTTP1()}
}

func (s *svc) Serve(ls net.Listener) error {
	server := &http.Server{
		Addr:    s.Address,
		Handler: s.buildRouter(),
	}
	return server.Serve(ls)
}

func (s *svc) buildRouter() http.Handler {
	// Server logger
	logger := log.New("golinks.http")

	// Create our HTTP Router.
	router := s.Gin

	// Configure HTTP Router Settings
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = false
	router.HandleMethodNotAllowed = false
	router.ForwardedByClientIP = true
	router.AppEngine = false
	router.UseRawPath = false
	router.UnescapePathValues = true

	// Log server config
	logger.Info("server link store: %s", s.LinkStore)
	logger.Info("server auth provider: %s", s.AuthProvider)

	// Setup middlewares.
	if s.Traced {
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
		Traced: s.Traced,
	})

	// Link module
	lnGroup := router.Group("links")
	lnGroup.Use(middlewares.Auth(s.AuthProvider))
	linkweb.Register(lnGroup, linkweb.Config{
		Store:  s.LinkStore,
		Traced: s.Traced,
	})

	// Org web module
	orgGroup := router.Group("org")
	authweb.Register(orgGroup, authweb.Config{
		Traced:   s.Traced,
		Provider: s.AuthProvider,
	})

	// Link api module
	lnAPIGroup := router.Group("api/links")
	lnAPIGroup.Use(middlewares.Auth(s.AuthProvider))
	linkapi.Register(lnAPIGroup, s.LinkStore)

	// Auth module
	authGroup := router.Group("api/orgs")
	authGroup.Use(middlewares.Auth(s.AuthProvider))
	authapi.Register(authGroup, s.AuthProvider)

	// Use redirect handler by default.
	router.NoRoute(
		middlewares.Auth(s.AuthProvider),
		middlewares.OrgRateLimit(5, time.Second),
		redirect.Handler(redirect.Config{
			Traced: s.Traced,
			Store:  s.LinkStore,
		}),
	)

	return router
}
