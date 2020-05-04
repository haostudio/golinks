package webbase

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"go.opencensus.io/trace"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/api/web"
)

// Base defines the web handler module.
type Base struct {
	web.Base
}

// Error defines the web error.
type Error struct {
	StatusCode int
	Messages   []string
	Log        string
}

// NewBase returns a new web handler module.
func NewBase(traced bool) Base {
	return Base{
		Base: web.NewBase(web.Config{
			PackrBox: packr.NewBox("../assets"),
			Traced:   traced,
		}),
	}
}

// Handler returns a gin handler rendering tmpl with data from dataGetter
// and 200 status code.
func (w *Base) Handler(
	tmpl string,
	dataGetter func(*gin.Context) (interface{}, *Error),
) gin.HandlerFunc {
	return w.handlerWithStatusCode(tmpl, http.StatusOK, dataGetter)
}

// UnAuthHandler returns a gin handler rendering tmpl with data from dataGetter
// and  401 status code.
func (w *Base) UnAuthHandler(tmpl string,
	dataGetter func(*gin.Context) (interface{}, *Error),
) gin.HandlerFunc {
	return w.handlerWithStatusCode(tmpl, http.StatusUnauthorized, dataGetter)
}

func (w *Base) handlerWithStatusCode(tmpl string, statusCode int,
	dataGetter func(*gin.Context) (interface{}, *Error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// traced
		if w.Traced {
			reqCtx, span := trace.StartSpan(ctx.Request.Context(), "web.handle")
			defer span.End()
			ctx.Request = ctx.Request.WithContext(reqCtx)
		}

		// handle
		logger := middlewares.GetLogger(ctx)
		pageData, err := dataGetter(ctx)
		if err != nil {
			logger.Error("failed to get page data. err: %v", err.Log)
			w.ServeErr(ctx, err)
			return
		}
		// render html
		if len(tmpl) > 0 && pageData != nil {
			w.Serve(ctx, statusCode, tmpl, pageData)
		}
	}
}

// ServeErr serves a html err page.
func (w *Base) ServeErr(ctx *gin.Context, err *Error) {
	if w.Traced {
		reqCtx, span := trace.StartSpan(ctx.Request.Context(), "web.serve_err")
		defer span.End()
		ctx.Request = ctx.Request.WithContext(reqCtx)

		span.AddAttributes(trace.StringAttribute("err", err.Log))
	}
	data := NewErrorPage()
	data.StatusCode = err.StatusCode
	data.ErrTitle = http.StatusText(err.StatusCode)
	data.Description = strings.Join(err.Messages, "; ")
	w.Serve(ctx, err.StatusCode, "error.html.tmpl", data)
}
