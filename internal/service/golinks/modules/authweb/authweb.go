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
			token, err := middlewares.GetToken(ctx)
			if err != nil {
				ctx.String(http.StatusOK, "logout success")
				return
			}
			err = conf.Manager.Logout(ctx.Request.Context(), token)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// remove token
			middlewares.DeleteToken(ctx)
			ctx.String(http.StatusOK, "logout success")
		})
	}

	// org/manage
	{
		manageRouter := router.Group("manage")
		manageRouter.Use(middlewares.AuthHTTPBasicAuth(conf.Manager))
		manageRouter.GET("", module.SetOrgUser())
		manageRouter.POST("", module.HandleSetOrgUserForm)
	}
}
