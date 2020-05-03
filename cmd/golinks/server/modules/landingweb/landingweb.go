package landingweb

import "github.com/gin-gonic/gin"

// Register register landing web in router.
func Register(router gin.IRouter, conf Config) {
	module := New(conf)
	router.GET("", module.Landing())
}
