package authweb

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

// Middleware returns the auth middleware that redirects to login page if the
// request is unauthorized.
func Middleware(manager *auth.Manager, path string) gin.HandlerFunc {
	return middlewares.Auth(manager, func(ctx *gin.Context) {
		callback := ctx.Request.URL.EscapedPath()
		url := fmt.Sprintf("/%s/login?callback=%s", strings.Trim(path, "/"), callback)
		ctx.Redirect(http.StatusMovedPermanently, url)
		ctx.Abort()
	})
}

// Config defines the web config.
type Config struct {
	Traced     bool
	Manager    *auth.Manager
	PathPrefix string
}

//Register register auth web in router
func Register(router gin.IRouter, conf Config) {
	web := New(conf)

	router.GET("login", func(ctx *gin.Context) {
		_, err := middlewares.GetUser(ctx)
		if err == nil {
			// already logged in
			ctx.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		if !errors.Is(err, middlewares.ErrNotFound) {
			web.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("failed to get user. %v", err),
			})
			return
		}
		web.Login()(ctx)
	})
	router.POST("login", web.HandleLoginForm)
	router.GET("logout", func(ctx *gin.Context) {
		token, err := middlewares.GetToken(ctx)
		if err != nil {
			web.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("failed to get token. %v", err),
			})
			return
		}
		err = conf.Manager.Logout(ctx.Request.Context(), token)
		if err != nil {
			web.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("logout failed. %v", err),
			})
			return
		}
		// remove token
		middlewares.DeleteToken(ctx)
		ctx.Redirect(http.StatusMovedPermanently, "/")
	})

	{
		orgRouter := router.Group("org")
		orgRouter.GET("register", web.SetOrg())
		orgRouter.POST("register", web.HandleSetOrgForm)

		orgRouter.Use(Middleware(conf.Manager, conf.PathPrefix))
		orgRouter.GET("manage", web.SetOrgUser())
		orgRouter.POST("manage", web.HandleSetOrgUserForm)
	}
}
