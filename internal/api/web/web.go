package web

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/gobuffalo/packr"
	"go.opencensus.io/trace"

	"github.com/haostudio/golinks/internal/api/middlewares"
)

// Config defines the web config struct.
type Config struct {
	PackrBox packr.Box
	Traced   bool
}

// Base defines the web handler module.
type Base struct {
	Config
}

// NewBase returns a new web handler module.
func NewBase(conf Config) Base {
	return Base{
		Config: conf,
	}
}

// Template returns a template instance that loads every tmpl file.
func (w *Base) Template() (t *template.Template, err error) {
	t = template.New("")
	err = w.PackrBox.Walk(func(name string, f packr.File) error {
		if name == "" {
			return nil
		}
		finfo, err := f.FileInfo()
		if err != nil {
			return err
		}
		// skip directory path
		if finfo.IsDir() {
			return nil
		}

		// skip all files end with .html.tmpl
		if !strings.HasSuffix(name, ".html.tmpl") {
			return nil
		}

		// Normalize template name
		n := name
		if strings.HasPrefix(name, "\\") || strings.HasPrefix(name, "/") {
			n = n[1:] // don't want template name to start with / ie. /index.html
		}
		// replace windows path separator \ to normalized /
		n = strings.Replace(n, "\\", "/", -1)

		str, err := w.PackrBox.FindString(name)
		if err != nil {
			return err
		}

		_, err = t.New(n).Parse(str)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t = nil
	}
	return
}

// Serve serves a tmpl file rendering with data.
func (w *Base) Serve(ctx *gin.Context,
	statusCode int, file string, data interface{}) {
	if w.Traced {
		reqCtx, span := trace.StartSpan(ctx.Request.Context(), "web.serve")
		defer span.End()
		ctx.Request = ctx.Request.WithContext(reqCtx)

		span.AddAttributes(trace.StringAttribute("serve_tmpl", file))
		span.AddAttributes(trace.Int64Attribute("serve_status", int64(statusCode)))
	}

	logger := middlewares.GetLogger(ctx)
	t, err := w.Template()
	if err != nil {
		logger.Error("failed to create template. err: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Render(statusCode, &render.HTML{
		Template: t,
		Name:     file,
		Data:     data,
	})
}
