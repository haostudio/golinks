package authweb

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
)

//Register register auth web in router
func Register(router gin.IRouter, conf Config) {
	module := New(conf)

	// org/
	{
		router.GET("", module.SetOrg())
		router.POST("", module.HandleSetOrgForm)
		// Logout
		router.GET("logout", func(ctx *gin.Context) {
			ctx.String(http.StatusUnauthorized, "logout success")
		})
	}

	// org/manage
	{
		manageRouter := router.Group("manage")
		manageRouter.Use(middlewares.Auth(conf.Manager))
		manageRouter.GET("", module.SetOrgUser())
		manageRouter.POST("", module.HandleSetOrgUserForm)
	}
}
