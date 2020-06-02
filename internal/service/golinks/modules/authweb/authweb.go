package authweb

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/service/golinks/ctx"
	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

// AuthRequired returns the middleware that requires logged in and redirects to
// login page if the request is unauthorized.
func AuthRequired(path string) gin.HandlerFunc {
	return ctx.AuthRequired(func(ginctx *gin.Context, err error) {
		if errors.Is(err, ctx.ErrNotFound) {
			callback := ginctx.Request.URL.EscapedPath()
			url := fmt.Sprintf(
				"/%s/login?callback=%s", strings.Trim(path, "/"), callback)
			ginctx.Redirect(http.StatusMovedPermanently, url)
			ginctx.Abort()
			return
		}
		ginctx.AbortWithStatus(http.StatusInternalServerError)
	})
}

// OrgRequired returns the middleware that requires organization and redirects
// to org register page if the user is without org.
func OrgRequired(path string) gin.HandlerFunc {
	return ctx.OrgRequired(func(ginctx *gin.Context, err error) {
		if errors.Is(err, ctx.ErrNotFound) {
			callback := ginctx.Request.URL.EscapedPath()
			url := fmt.Sprintf(
				"/%s/org/register?callback=%s", strings.Trim(path, "/"), callback)
			ginctx.Redirect(http.StatusMovedPermanently, url)
			ginctx.Abort()
			return
		}
		ginctx.AbortWithStatus(http.StatusInternalServerError)
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

	router.GET("login", func(ginctx *gin.Context) {
		_, err := ctx.GetUser(ginctx)
		if err == nil {
			// already logged in
			ginctx.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		if !errors.Is(err, ctx.ErrNotFound) {
			web.ServeErr(ginctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("failed to get user. %v", err),
			})
			return
		}
		web.Login()(ginctx)
	})
	router.POST("login", web.HandleLoginForm)
	router.GET("logout", func(ginctx *gin.Context) {
		token, err := ctx.GetToken(ginctx)
		if err != nil {
			web.ServeErr(ginctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("failed to get token. %v", err),
			})
			return
		}
		err = conf.Manager.Logout(ginctx.Request.Context(), token)
		if err != nil {
			web.ServeErr(ginctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("logout failed. %v", err),
			})
			return
		}
		// remove token
		ctx.DeleteToken(ginctx)
		ginctx.Redirect(http.StatusMovedPermanently, "/")
	})

	{
		orgRouter := router.Group("org")
		orgRouter.Use(AuthRequired(conf.PathPrefix))
		orgRouter.GET("register", func(ginctx *gin.Context) {
			_, err := ctx.GetOrg(ginctx)
			if err == nil {
				// already in an org
				ginctx.Redirect(
					http.StatusMovedPermanently,
					fmt.Sprintf("/%s/org/manage", conf.PathPrefix),
				)
				return
			}
			if !errors.Is(err, ctx.ErrNotFound) {
				web.ServeErr(ginctx, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
					Log:        fmt.Sprintf("failed to get org. %v", err),
				})
				return
			}
			web.OrgRegister()(ginctx)
		})
		orgRouter.POST("register", web.HandleOrgRegisterForm)

		orgRouter.Use(OrgRequired(conf.PathPrefix))
		orgRouter.GET("manage", web.SetOrgUser())
		orgRouter.POST("manage", web.HandleSetOrgUserForm)
	}
}
