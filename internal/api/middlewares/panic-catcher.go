package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// PanicCatcher defines a panic catcher handler.
func PanicCatcher(ctx *gin.Context) {
	// Recover From Panic & Log Stack Trace to StackDriver
	defer func() {
		if recovered := recover(); recovered != nil {
			logger := GetLogger(ctx)
			logger.Error("\x1b[31m%v\n[Stack Trace]\n%s\x1b[m",
				recovered, debug.Stack())
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	// Process Request Chain
	ctx.Next()
}
