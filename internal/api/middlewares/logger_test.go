package middlewares

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/popodidi/log"
)

func TestMiddlewareFirst(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{
		URL: &url.URL{},
	}

	logger := GetLogger(ctx)
	require.NotNil(t, logger)

	_, logLogger := log.Context(ctx.Request.Context(), "test")
	require.NotNil(t, logLogger)

	require.Equal(t, logLogger.GetID(), logger.GetID())
}

func TestLoggerFirst(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{
		URL: &url.URL{},
	}

	reqCtx, logLogger := log.Context(ctx.Request.Context(), "test")
	ctx.Request = ctx.Request.WithContext(reqCtx)
	require.NotNil(t, logLogger)

	logger := GetLogger(ctx)
	require.NotNil(t, logger)

	require.Equal(t, logLogger.GetID(), logger.GetID())
}
