package authapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/auth"
)

// New returns auth api handles.
func New(manager *auth.Manager) *Auth {
	return &Auth{
		manager: manager,
	}
}

// Auth defines the auth api handlers struct.
type Auth struct {
	manager *auth.Manager
}

// PathParamUserKey returns the user path parameter
func (a *Auth) PathParamUserKey() string {
	return "user"
}

// PathParamOrgKey returns the org path parameter
func (a *Auth) PathParamOrgKey() string {
	return "org"
}

// SetOrg sets org.
func (a *Auth) SetOrg(ctx *gin.Context) {
	logger := middlewares.GetLogger(ctx)
	key := ctx.Param(a.PathParamOrgKey())
	if len(key) == 0 {
		logger.Error("Auth: SetOrg: Empty key")
		ctx.String(http.StatusBadRequest, "empty key")
		return
	}

	// check if the org exists
	_, err := a.manager.GetOrg(ctx, key)
	if !errors.Is(err, auth.ErrNotFound) {
		logger.Error("org already exists, err")
		ctx.Status(http.StatusBadRequest)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err = ctx.BindJSON(&req)
	if err != nil {
		logger.Error("failed to bind json, err: %v", err)
		ctx.String(http.StatusBadRequest, "parameters error")
		return
	}
	org := auth.Organization{
		Name:       key,
		AdminEmail: req.Email,
	}
	err = a.manager.SetOrg(ctx.Request.Context(), org)
	if err != nil {
		logger.Error("failed to set org, err: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	user, err := auth.NewUser(req.Email, req.Password, key)
	if err != nil {
		logger.Error("failed to set user pwd, err: %v", err)
		err = a.manager.DeleteOrg(ctx.Request.Context(), key)
		if err != nil {
			logger.Error("failed to del org, err: %v", err)
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	err = a.manager.SetUser(ctx.Request.Context(), *user)
	if err != nil {
		logger.Error("failed to set user, err: %v", err)
		err = a.manager.DeleteOrg(ctx.Request.Context(), key)
		if err != nil {
			logger.Error("failed to del org, err: %v", err)
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	res := make(map[string]string)
	res["Organization"] = org.Name
	res["AdminEmail"] = org.AdminEmail

	ctx.JSON(http.StatusOK, org)
}

// GetOrg returns org.
func (a *Auth) GetOrg(ctx *gin.Context) {
	logger := middlewares.GetLogger(ctx)
	key := ctx.Param(a.PathParamOrgKey())
	if len(key) == 0 {
		logger.Error("Auth: GetOrg: Empty key")
		ctx.String(http.StatusBadRequest, "empty key")
		return
	}
	org, err := a.manager.GetOrg(ctx.Request.Context(), key)
	if err != nil {
		logger.Error("failed to get org, err: %v", err)
		ctx.Status(http.StatusInternalServerError)
	}
	ctx.JSON(http.StatusOK, org)
}

// GetOrgUsers returns org users.
func (a *Auth) GetOrgUsers(ctx *gin.Context) {
	logger := middlewares.GetLogger(ctx)
	key := ctx.Param(a.PathParamOrgKey())
	if len(key) == 0 {
		logger.Error("Auth: GetOrg: Empty key")
		ctx.String(http.StatusBadRequest, "empty key")
		return
	}
	users, err := a.manager.GetOrgUsers(ctx.Request.Context(), key)
	if err != nil {
		logger.Error("failed to get org users, err: %v", err)
		ctx.Status(http.StatusInternalServerError)
	}
	ctx.JSON(http.StatusOK, users)
}

// GetUser returns a user.
func (a *Auth) GetUser(ctx *gin.Context) {
	logger := middlewares.GetLogger(ctx)
	key := ctx.Param(a.PathParamUserKey())
	if len(key) == 0 {
		logger.Error("Auth: GetUser: Empty key")
		ctx.String(http.StatusBadRequest, "empty key")
		return
	}

	user, err := a.manager.GetUser(ctx.Request.Context(), key)
	if errors.Is(err, auth.ErrNotFound) {
		logger.Error("User Not Found")
		ctx.String(http.StatusBadRequest, "user not found")
		return
	}
	if err != nil {
		logger.Error("Get Users Error:%v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	res := make(map[string]string)
	res["Email"] = user.Email
	res["Organization"] = user.Organization

	ctx.JSON(http.StatusOK, res)
}

// GetUsers returns all users.
func (a *Auth) GetUsers(ctx *gin.Context) {
	logger := middlewares.GetLogger(ctx)
	users, err := a.manager.GetUsers(ctx.Request.Context())
	if errors.Is(err, auth.ErrNotFound) {
		logger.Error("User Not Found")
		ctx.String(http.StatusBadRequest, "user not found")
		return
	}
	if err != nil {
		logger.Error("Get Users Error:%v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// SetOrgUser sets org user.
func (a *Auth) SetOrgUser(ctx *gin.Context) {
	logger := middlewares.GetLogger(ctx)
	key := ctx.Param(a.PathParamOrgKey())
	if len(key) == 0 {
		logger.Error("Auth: SetOrgUser: Empty key")
		ctx.String(http.StatusBadRequest, "empty key")
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ctx.BindJSON(&req)
	if err != nil {
		logger.Error("failed to bind json, err: %v", err)
		ctx.String(http.StatusBadRequest, "parameters error")
		return
	}

	user, err := auth.NewUser(req.Email, req.Password, key)
	if err != nil {
		logger.Error("failed to init user, err: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	err = a.manager.SetUser(ctx.Request.Context(), *user)
	if err != nil {
		logger.Error("failed to set user, err: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	res := make(map[string]string)
	res["organization"] = user.Organization
	res["email"] = user.Email

	ctx.JSON(http.StatusOK, res)
}
