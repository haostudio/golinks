package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/auth"
)

// NoAuth returns a no auth handler with default org.
func NoAuth(defaultOrg string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(orgKey, auth.Organization{
			Name: defaultOrg,
		})
	}
}

// AuthRequired returns the auth required middleware based on the
// "GOLINKS_TOKEN" cookie.
func AuthRequired(onError func(*gin.Context, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := GetLogger(ctx)
		_, err := GetUser(ctx)
		if err != nil {
			logger.Error("failed to get user. err: %v", err)
			onError(ctx, err)
			return
		}
	}
}

// AuthSimple401 returns the auth required middleware based on the
// "GOLINKS_TOKEN" cookie and returns 401 if the user is unauthorized.
var AuthSimple401 = AuthRequired(func(ctx *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.AbortWithStatus(http.StatusInternalServerError)
})

// OrgRequired returns the org required middleware based on the
// "GOLINKS_TOKEN" cookie.
func OrgRequired(onError func(*gin.Context, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := GetLogger(ctx)
		_, err := GetOrg(ctx)
		if err != nil {
			logger.Error("failed to get org. err: %v", err)
			onError(ctx, err)
			return
		}
	}
}

// OrgSimple404 returns the org required middleware based on the
// "GOLINKS_TOKEN" cookie and returns 404 if the user is unauthorized.
var OrgSimple404 = OrgRequired(func(ctx *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.AbortWithStatus(http.StatusInternalServerError)
})

// GetToken returns the token cookie.
func GetToken(ctx *gin.Context) (token string, err error) {
	GetLogger(ctx).Debug("get token")
	return ctx.Cookie(tokenCookieKey)
}

// SetToken sets the token cookie
func SetToken(ctx *gin.Context, token string, maxAge int) {
	GetLogger(ctx).Debug("set token")
	ctx.SetCookie(tokenCookieKey, token, maxAge, "", "", false, false)
}

// DeleteToken deletes the token cookie
func DeleteToken(ctx *gin.Context) {
	GetLogger(ctx).Debug("delete token")
	ctx.SetCookie(tokenCookieKey, "", 0, "", "", false, false)
}
