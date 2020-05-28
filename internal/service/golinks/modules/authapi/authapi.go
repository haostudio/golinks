package authapi

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/auth"
)

// Register registers auth endpoints in router.
func Register(router gin.IRouter, manager *auth.Manager) {
	module := New(manager)
	router.GET(
		fmt.Sprintf("/:%s", module.PathParamOrgKey()),
		module.GetOrg,
	)
	router.POST(
		fmt.Sprintf("/:%s", module.PathParamOrgKey()),
		module.SetOrg,
	)
	router.POST(
		fmt.Sprintf("/:%s/user", module.PathParamOrgKey()),
		module.SetOrgUser,
	)
	router.GET(
		fmt.Sprintf("/:%s/users", module.PathParamOrgKey()),
		module.GetOrgUsers,
	)
}
