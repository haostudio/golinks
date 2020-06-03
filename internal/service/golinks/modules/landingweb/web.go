package landingweb

import (
	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

// Config defines the web config.
type Config struct {
	Traced bool
}

// Web defines the web handler module.
type Web struct {
	webbase.Base
}

// New returns a new web handler module.
func New(conf Config) *Web {
	return &Web{
		Base: webbase.NewBase(conf.Traced),
	}
}

// Landing returns the landing page. (./web/landing.html)
func (w *Web) Landing() gin.HandlerFunc {
	return w.Handler(
		"landing.html.tmpl",
		func(ctx *gin.Context) (interface{}, *webbase.Error) {
			return NewData(ctx), nil
		},
	)
}
