package middlewares

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/auth"
)

// Auth defines the auth middleware.
func Auth(provider auth.Provider) gin.HandlerFunc {
	// use HTTP basic auth
	realm := "Basic realm=" + strconv.Quote("Authorization Required")
	onAuthError := func(ctx *gin.Context) {
		ctx.Header("WWW-Authenticate", realm)
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	return func(ctx *gin.Context) {
		logger := GetLogger(ctx)
		authHeader := ctx.Request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			onAuthError(ctx)
			return
		}
		if len(authHeader) <= len("Basic ") {
			logger.Debug("invalid header: %s", authHeader)
			onAuthError(ctx)
			return
		}
		// decode base64
		authHeader = authHeader[len("Basic "):]
		authStr, err := base64.StdEncoding.DecodeString(authHeader)
		if err != nil {
			logger.Debug("invalid header: %s", authHeader)
			onAuthError(ctx)
			return
		}
		authStrs := strings.Split(string(authStr), ":")
		if len(authStrs) != 2 {
			logger.Debug("invalid auth string: %s", authStr)
			onAuthError(ctx)
			return
		}
		user, err := provider.GetUser(ctx.Request.Context(), authStrs[0])
		if errors.Is(err, auth.ErrNotFound) {
			logger.Debug("invalid user email: %s", authStrs[0])
			onAuthError(ctx)
			return
		} else if err != nil {
			logger.Error("failed to get user. err: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if err := user.VerifyPassword(authStrs[1]); err != nil {
			logger.Debug("invalid user password of %s. err: %v", authStrs[0], err)
			onAuthError(ctx)
			return
		}

		// Authorized
		getOrg := func(ctx *gin.Context) (org auth.Organization, err error) {
			return provider.GetOrg(ctx.Request.Context(), user.Organization)
		}
		ctx.Set(userKey, user)
		ctx.Set(getOrgKey, getOrg)
	}
}

// GetUser returns the user stored in gin.Context.
func GetUser(ctx *gin.Context) (user auth.User, err error) {
	u, ok := ctx.Get(userKey)
	if !ok {
		err = ErrNotFound
		return
	}
	user, ok = u.(auth.User)
	if !ok {
		err = ErrInternal
		return
	}
	return
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
