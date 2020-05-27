package golinks

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
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
	Gin     *gin.Engine
	Address string
	Traced  bool
	Wiki    bool // server static doc site
	Auth    struct {
		Enabled    bool
		DefaultOrg string        // default org for Auth.Enabled = false
		Manager    *auth.Manager // provider for Auth.Enabled = true
	}
	LinkStore link.Store
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
	if s.Auth.Enabled {
		logger.Info("server auth provider: %s", s.Auth.Manager)
	} else {
		logger.Warn("server auth disabled")
		logger.Warn("server default org: %s", s.Auth.DefaultOrg)
	}

	// Setup middlewares.
	if s.Traced {
		router.Use(middlewares.Trace)
	}
	router.Use(middlewares.CtxLogger)
	router.Use(middlewares.PanicCatcher)

	// static doc site
	if s.Wiki {
		fileServer := http.FileServer(packr.NewBox("./wiki"))
		router.Any("wiki/*path", func(ctx *gin.Context) {
			// Trim the wiki prefix
			ctx.Request.URL.Path = ctx.Param("path")
			fileServer.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	// nolint: godox
	// FIXME: add favicon and remove this hack
	router.Use(func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.URL.String(), "/favicon.ico") {
			ctx.AbortWithStatus(http.StatusNotFound)
		}
	})

	// Landing page
	landingweb.Register(router, landingweb.Config{
		Traced:      s.Traced,
		AuthEnabled: s.Auth.Enabled,
	})

	// Link web module
	lnGroup := router.Group("links")
	if s.Auth.Enabled {
		lnGroup.Use(middlewares.AuthHTTPBasicAuth(s.Auth.Manager))
	} else {
		lnGroup.Use(middlewares.NoAuth(s.Auth.DefaultOrg))
	}
	linkweb.Register(lnGroup, linkweb.Config{
		Store:  s.LinkStore,
		Traced: s.Traced,
	})

	// Org web module
	if s.Auth.Enabled {
		orgGroup := router.Group("org")
		authweb.Register(orgGroup, authweb.Config{
			Traced:  s.Traced,
			Manager: s.Auth.Manager,
		})
	}

	// Link api module
	lnAPIGroup := router.Group("api/links")
	if s.Auth.Enabled {
		lnAPIGroup.Use(middlewares.AuthSimple401(s.Auth.Manager))
	} else {
		lnAPIGroup.Use(middlewares.NoAuth(s.Auth.DefaultOrg))
	}
	linkapi.Register(lnAPIGroup, s.LinkStore)

	// Auth module
	if s.Auth.Enabled {
		authGroup := router.Group("api/orgs")
		authGroup.Use(middlewares.AuthHTTPBasicAuth(s.Auth.Manager))
		authapi.Register(authGroup, s.Auth.Manager)
	}

	// nolint: godox
	// TODO: configure rate limit from golinks.Config
	// Use redirect handler by default.
	var noRoute []gin.HandlerFunc
	if s.Auth.Enabled {
		noRoute = append(noRoute, middlewares.AuthHTTPBasicAuth(s.Auth.Manager))
		noRoute = append(noRoute, middlewares.OrgRateLimit(5, time.Second))
	} else {
		noRoute = append(noRoute, middlewares.NoAuth(s.Auth.DefaultOrg))
		noRoute = append(noRoute, middlewares.OrgRateLimit(100, time.Second))
	}
	noRoute = append(noRoute,
		redirect.Handler(redirect.Config{
			Traced: s.Traced,
			Store:  s.LinkStore,
		}),
	)
	router.NoRoute(noRoute...)

	return router
}
