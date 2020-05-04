package linkapi

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/link"
)

// Register register api in router.
func Register(router gin.IRouter, lnStore link.Store) {
	module := New(lnStore)
	router.GET("", module.GetLinks)
	// Admin functions
	router.PUT(
		fmt.Sprintf(":%s", module.PathParamLinkKey()),
		module.UpdateLink,
	)
	router.DELETE(
		fmt.Sprintf(":%s", module.PathParamLinkKey()),
		module.DeleteLink,
	)
}
