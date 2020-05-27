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
		tokenStr, err := GetToken(ctx)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				onAuthError(ctx)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		claims, err := manager.Verify(ctx.Request.Context(), tokenStr)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) ||
				errors.Is(err, auth.ErrTokenExpired) {
				onAuthError(ctx)
				return
			}
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

// AuthHTTPBasicAuth returns the auth middleware based on the "GOLINKS_TOKEN"
// cookie and requires HTTP Basic Authentication if the user is unauthorized.
func AuthHTTPBasicAuth(manager *auth.Manager) gin.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote("Authorization Required")
	onAuthError := func(ctx *gin.Context) {
		ctx.Header("WWW-Authenticate", realm)
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	return Auth(manager, func(ctx *gin.Context) {
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
		email, password := authStrs[0], authStrs[1]
		token, err := manager.Login(ctx.Request.Context(), email, password)
		if err != nil {
			logger.Debug("manager login error. %s. %s", email)
			if errors.Is(err, ErrNotFound) {
				onAuthError(ctx)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		claims, err := manager.Verify(ctx.Request.Context(), token.JWT)
		if err != nil {
			logger.Debug("manager verify error. s. %s", email, err)
			if errors.Is(err, ErrNotFound) {
				onAuthError(ctx)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// authorized
		SetToken(ctx, token.JWT, int(manager.TokenExpieration.Seconds()))
		getUser := func(ctx *gin.Context) (user auth.User, err error) {
			return manager.GetUser(ctx.Request.Context(), email)
		}
		getOrg := func(ctx *gin.Context) (org auth.Organization, err error) {
			return manager.GetOrg(ctx.Request.Context(), claims.Org)
		}
		ctx.Set(getUserKey, getUser)
		ctx.Set(getOrgKey, getOrg)
	})
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
	return ctx.Cookie(tokenCookieKey)
}

// SetToken sets the token cookie
func SetToken(ctx *gin.Context, token string, maxAge int) {
	ctx.SetCookie(
		tokenCookieKey,
		token,
		maxAge,
		"", "",
		true, true,
	)
}

// DeleteToken deletes the token cookie
func DeleteToken(ctx *gin.Context) {
	ctx.SetCookie(tokenCookieKey, "", 0, "", "", true, true)
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
