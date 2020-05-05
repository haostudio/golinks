package landingweb

import (
	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

// Config defines the web config.
type Config struct {
	AuthEnabled bool
	Traced      bool
}

// Web defines the web handler module.
type Web struct {
	webbase.Base
	AuthEnabled bool
}

// New returns a new web handler module.
func New(conf Config) *Web {
	return &Web{
		Base:        webbase.NewBase(conf.Traced),
		AuthEnabled: conf.AuthEnabled,
	}
}

// Landing returns the landing page. (./web/landing.html)
func (w *Web) Landing() gin.HandlerFunc {
	return w.Handler(
		"landing.html.tmpl",
		func(ctx *gin.Context) (interface{}, *webbase.Error) {
			return NewData(w.AuthEnabled), nil
		},
	)
}
