package redirect

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/link"
	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

// Config defines the config struct.
type Config struct {
	Traced bool
	Store  link.Store
}

// Handler redirects requests based on the link.Store.
func Handler(conf Config) gin.HandlerFunc {
	web := webbase.NewBase(conf.Traced)
	return func(ctx *gin.Context) {
		logger := middlewares.GetLogger(ctx)

		path := ctx.Request.URL.Path
		key, param := link.Parse(path)
		logger.Debug("key=%s param=%s", key, param)

		org, err := middlewares.GetOrg(ctx)
		if err != nil {
			logger.Error("failed to get org. err: %v", err)
			web.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
			})
			return
		}
		ln, err := conf.Store.GetLink(ctx, org.Name, key)
		if errors.Is(err, link.ErrNotFound) {
			ctx.Redirect(
				http.StatusTemporaryRedirect, fmt.Sprintf("/links/edit/%s", key))
			return
		}
		if err != nil {
			logger.Error("failed to get link from store. err: %v", err)
			web.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
			})
			return
		}
		// Link Found!
		target, err := ln.GetRedirectLink(param)
		if errors.Is(err, link.ErrInvalidParams) {
			logger.Error("invalid param")
			desc, err := ln.Description()
			if err != nil {
				logger.Error("failed to get link desc. err: %v", err)
				web.ServeErr(ctx, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
				})
				return
			}
			web.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Invalid params", desc},
			})
			return
		} else if err != nil {
			logger.Error("failed to get target link. err: %v", err)
			web.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
			})
			return
		}
		logger.Debug("redirect %s to %s", path, target)
		ctx.Redirect(http.StatusTemporaryRedirect, target)
	}
}
