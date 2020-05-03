package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/popodidi/log"
)

// Key defines the key for logger to save in gin.Context.
const (
	MethodTag = "golinks.middlewares.logger.method"
	URLTag    = "golinks.middlewares.logger.url"
)

// GetLogger gets logger from gin context and returns log.Null() as the
// default value.
func GetLogger(c *gin.Context) log.Logger {
	logger := log.GetFromCtx(c.Request.Context())
	if logger != nil {
		return logger
	}

	var ctx context.Context
	ctx, logger = log.Context(c.Request.Context())
	logger.GetLabels().Set(MethodTag, c.Request.Method)
	logger.GetLabels().Set(URLTag, c.Request.URL.String())
	c.Request = c.Request.WithContext(ctx)
	return logger
}
