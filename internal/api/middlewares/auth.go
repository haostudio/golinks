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
		getOrg := func(ctx *gin.Context) (org auth.Organization, err error) {
			org.Name = defaultOrg
			org.AdminEmail = ""
			return
		}
		ctx.Set(getOrgKey, getOrg)
	}
}

// Auth returns the auth middleware based on the "GOLINKS_TOKEN"
// cookie.
func Auth(manager *auth.Manager, onAuthError gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := GetLogger(ctx)
		tokenStr, err := GetToken(ctx)
		if errors.Is(err, http.ErrNoCookie) || tokenStr == "" {
			logger.Error("cookie not found")
			onAuthError(ctx)
			return
		}
		if err != nil {
			logger.Error("failed to get cookie. err: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		claims, err := manager.Verify(ctx.Request.Context(), tokenStr)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) ||
				errors.Is(err, auth.ErrTokenExpired) {
				logger.Error("invalid token. err: %v", err)
				onAuthError(ctx)
				return
			}
			logger.Error("failed to verify token. err: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// authorized
		getUser := func(ctx *gin.Context) (user auth.User, err error) {
			return manager.GetUser(ctx.Request.Context(), claims.Email)
		}
		getOrg := func(ctx *gin.Context) (org auth.Organization, err error) {
			return manager.GetOrg(ctx.Request.Context(), claims.Org)
		}
		ctx.Set(getUserKey, getUser)
		ctx.Set(getOrgKey, getOrg)
	}
}

// AuthSimple401 returns the auth middleware based on the "GOLINKS_TOKEN"
// cookie and returns 401 if the user is unauthorized.
func AuthSimple401(manager *auth.Manager) gin.HandlerFunc {
	return Auth(manager, func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	})
}

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

// GetUser returns the user stored in gin.Context.
func GetUser(ctx *gin.Context) (user auth.User, err error) {
	val, ok := ctx.Get(getUserKey)
	if !ok {
		err = ErrNotFound
		return
	}
	f, ok := val.(func(*gin.Context) (auth.User, error))
	if !ok {
		err = ErrInternal
		return
	}
	return f(ctx)
}

// GetOrg returns the org stored in gin.Context.
func GetOrg(ctx *gin.Context) (org auth.Organization, err error) {
	val, ok := ctx.Get(getOrgKey)
	if !ok {
		err = ErrNotFound
		return
	}
	f, ok := val.(func(*gin.Context) (auth.Organization, error))
	if !ok {
		err = ErrInternal
		return
	}
	return f(ctx)
}
