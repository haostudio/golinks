package linkweb

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Register register links web in router.
func Register(router gin.IRouter, conf Config) {
	module := New(conf)
	router.GET("", module.Links())
	// Admin pages
	router.GET(
		fmt.Sprintf("edit/:%s", module.PathParamLinkKey()),
		module.EditLink(),
	)
	router.POST(
		fmt.Sprintf("edit/:%s", module.PathParamLinkKey()),
		module.HandleEditLinktForm,
	)
}
