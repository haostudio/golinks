package ctx

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	slidingwindow "github.com/haostudio/golinks/internal/limiter/sliding-window"
)

// OrgRateLimit defines a rate limited middleware of an org.
func OrgRateLimit(limit int, windowSize time.Duration) gin.HandlerFunc {
	return middlewares.Limited(
		func(ctx *gin.Context) (string, error) {
			logger := middlewares.GetLogger(ctx)
			org, err := GetOrg(ctx)
			if err != nil {
				logger.Error("failed to get org. err: %v", err)
				return "", err
			}
			return org.Name, nil
		},
		slidingwindow.New(limit, windowSize),
	)
}
