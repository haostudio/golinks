package ctx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/auth"
)

// Ctx defines the golinks service context.
type Ctx struct {
	Auth struct {
		Enabled bool
		Manager *auth.Manager
	}
}

// Middeware prepares golinks service context in gin.Context.
func Middeware(ctx Ctx) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ginctx.Set(ctxKey, ctx)
	}
}

// IsAuthEnabled returns the auth is enabled.
func IsAuthEnabled(ctx *gin.Context) bool {
	val, ok := ctx.Get(ctxKey)
	if !ok {
		return false
	}
	return val.(Ctx).Auth.Enabled
}

// GetAuthManager returns the auth.Manager set in context.
func GetAuthManager(ctx *gin.Context) (manager *auth.Manager, ok bool) {
	val, ok := ctx.Get(ctxKey)
	if !ok {
		return
	}
	manager = val.(Ctx).Auth.Manager
	if manager == nil {
		ok = false
	}
	return
}

// GetUser returns the user of the request.
func GetUser(ctx *gin.Context) (user auth.User, err error) {
	// get cached value
	val, ok := ctx.Get(userKey)
	if ok {
		user = val.(auth.User)
		return
	}
	// get from manager
	logger := middlewares.GetLogger(ctx)
	manager, ok := GetAuthManager(ctx)
	if !ok {
		logger.Error("auth manager not found")
		err = ErrNotFound
		return
	}
	// get from token
	tokenStr, err := GetToken(ctx)
	if errors.Is(err, http.ErrNoCookie) || tokenStr == "" {
		logger.Error("cookie not found")
		err = ErrNotFound
		return
	}
	if err != nil {
		logger.Error("failed to get cookie. err: %v", err)
		return
	}
	claims, err := manager.Verify(ctx.Request.Context(), tokenStr)
	if err != nil {
		logger.Error("failed to verify token. err: %v", err)
		if errors.Is(err, auth.ErrInvalidToken) ||
			errors.Is(err, auth.ErrTokenExpired) {
			err = ErrNotFound
			return
		}
		err = ErrInternal
		return
	}
	user, err = manager.GetUser(ctx.Request.Context(), claims.Email)
	if err != nil {
		return
	}
	ctx.Set(userKey, user)
	return
}

// GetOrg returns the org of the request.
func GetOrg(ctx *gin.Context) (org auth.Organization, err error) {
	// get cached value
	val, ok := ctx.Get(orgKey)
	if ok {
		org = val.(auth.Organization)
		return
	}
	logger := middlewares.GetLogger(ctx)
	user, err := GetUser(ctx)
	if err != nil {
		logger.Error("failed to get user. err: %v", err)
		return
	}
	if user.Organization == "" {
		err = ErrNotFound
		return
	}
	manager, ok := GetAuthManager(ctx)
	if !ok {
		logger.Error("auth manager not found")
		err = ErrNotFound
		return
	}
	org, err = manager.GetOrg(ctx.Request.Context(), user.Organization)
	if err != nil {
		logger.Error("failed to get org. err: %v", err)
		return
	}
	ctx.Set(orgKey, org)
	return
}
