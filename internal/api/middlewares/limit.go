package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/limiter"
)

// Limited returns a middleware that limits the access of clients with
// keyFunc and limiter.
func Limited(
	keyFunc func(*gin.Context) (string, error),
	limit limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := GetLogger(ctx)
		rateLimitKey, err := keyFunc(ctx)
		if err != nil {
			logger.Error("failed to get rate limit key. err: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		res, err := limit.TryHit(ctx, rateLimitKey)
		if err != nil {
			logger.Error("failed to get rate limit result. err: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(res.Limit))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(res.Reset, 10))

		if res.Reached {
			ctx.Status(http.StatusTooManyRequests)
			return
		}
	}
}
