package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// CtxLogger logs the request context
func CtxLogger(ctx *gin.Context) {
	logger := GetLogger(ctx)
	address := ctx.ClientIP()
	method := ctx.Request.Method
	url := ctx.Request.URL.EscapedPath()
	headerCopy := make(map[string][]string)
	for k, v := range ctx.Request.Header {
		headerCopy[k] = make([]string, len(v))
		copy(headerCopy[k], v)
	}
	headersStr := fmt.Sprintf("%v", headerCopy)

	// execute handlers.
	ctx.Next()

	// get status.
	statusCode := ctx.Writer.Status()
	desc := fmt.Sprintf("%s | %3d %-6s %s", address, statusCode, method, url)
	if statusCode >= 200 && statusCode < 300 {
		logger.Info(desc)
	} else {
		logger.Warn(desc)
		logger.Debug("headers: %s", headersStr)
	}
}
