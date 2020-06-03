package authweb

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/service/golinks/ctx"
	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

const (
	formBtnActionCreate   = "create"
	formBtnActionLogin    = "login"
	formBtnActionRegister = "register"
	formBtnActionSave     = "save"

	formInputAction   = "action"
	formInputCallback = "callback"
	formInputEmail    = "email"
	formInputName     = "name"
	formInputPassword = "password"
)

// Web defines the web handler module.
type Web struct {
	webbase.Base
	manager *auth.Manager
}

// New returns a new web handler module.
func New(conf Config) *Web {
	return &Web{
		Base:    webbase.NewBase(conf.Traced),
		manager: conf.Manager,
	}
}

// Login returns the login and register page.
func (w *Web) Login() gin.HandlerFunc {
	return w.Handler(
		"login.html.tmpl",
		func(ginctx *gin.Context) (interface{}, *webbase.Error) {
			callback := ginctx.Query("callback")
			if callback == "" {
				callback = "/"
			}
			return LoginData{
				Data:                  webbase.NewData("Golinks - Login", ginctx),
				FormInputEmail:        formInputEmail,
				FormInputPassword:     formInputPassword,
				FormInputCallback:     formInputCallback,
				FormInputAction:       formInputAction,
				FormLoginBtnAction:    formBtnActionLogin,
				FormRegisterBtnAction: formBtnActionRegister,
				Callback:              callback,
			}, nil
		},
	)
}

// HandleLoginForm handle request to create org user.
func (w *Web) HandleLoginForm(ginctx *gin.Context) {
	email := ginctx.PostForm(formInputEmail)
	password := ginctx.PostForm(formInputPassword)
	action := ginctx.PostForm(formInputAction)
	callback := ginctx.PostForm(formInputCallback)
	switch action {
	case formBtnActionLogin:
		token, err := w.manager.Login(ginctx.Request.Context(), email, password)
		if err != nil {
			w.ServeErr(ginctx, &webbase.Error{
				StatusCode: http.StatusUnauthorized,
				Messages:   []string{"Invalid email/password"},
				Log:        fmt.Sprintf("login failed. %v", err),
			})
			return
		}
		// authorized
		ctx.SetToken(
			ginctx, token.JWT, int(w.manager.TokenExpieration.Seconds()))
		ginctx.Redirect(http.StatusMovedPermanently, callback)
		return
	case formBtnActionRegister:
		user, err := auth.NewUser(email, password, "")
		if err != nil {
			w.ServeErr(ginctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("failed to create user. %v", err),
			})
		}
		err = w.manager.RegisterUser(ginctx.Request.Context(), *user)
		if errors.Is(err, auth.ErrUserExists) {
			w.ServeErr(ginctx, &webbase.Error{
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Already registered"},
				Log:        fmt.Sprintf("failed to create user. %v", err),
			})
			return
		}
		if err != nil {
			w.ServeErr(ginctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("failed to register user. %v", err),
			})
			return
		}
		token, err := w.manager.Login(ginctx.Request.Context(), email, password)
		if err != nil {
			w.ServeErr(ginctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log:        fmt.Sprintf("login failed. %v", err),
			})
			return
		}
		// authorized
		ctx.SetToken(
			ginctx, token.JWT, int(w.manager.TokenExpieration.Seconds()))
		ginctx.Redirect(http.StatusMovedPermanently, "/")
		return
	default:
		w.ServeErr(ginctx, &webbase.Error{
			StatusCode: http.StatusBadRequest,
			Messages:   []string{"Invalid action"},
			Log:        fmt.Sprintf("invalid action %s", action),
		})
		return
	}
}

// SetOrgUser sets org user.
func (w *Web) SetOrgUser() gin.HandlerFunc {
	return w.Handler(
		"auth.html.tmpl",
		func(ginctx *gin.Context) (interface{}, *webbase.Error) {
			org, err := ctx.GetOrg(ginctx)
			if err != nil {
				return nil, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
					Log:        fmt.Sprintf("failed to get org, err: %v", err),
				}
			}

			user, err := ctx.GetUser(ginctx)
			if err != nil {
				return nil, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
					Log:        fmt.Sprintf("failed to get user, err: %v", err),
				}
			}
			admin := false
			if org.AdminEmail == user.Email {
				admin = true
			}

			users, err := w.manager.GetOrgUsers(ginctx, org.Name)
			if err != nil {
				return nil, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
					Log:        fmt.Sprintf("failed to get org users, err: %v", err),
				}
			}

			userSlice := make([]User, len(users))
			for i, d := range users {
				userSlice[i] = User{Email: d}
			}

			return PageData{
				Data: webbase.NewData(
					"Golinks - Organization Management", ginctx),
				Users:          userSlice,
				Admin:          admin,
				FormInputEmail: formInputEmail,
				FormBtnAction:  formBtnActionSave,
			}, nil
		},
	)
}

// HandleSetOrgUserForm handle request to create org user.
func (w *Web) HandleSetOrgUserForm(ginctx *gin.Context) {
	email := ginctx.PostForm(formInputEmail)
	org, err := ctx.GetOrg(ginctx)
	if err != nil {
		w.ServeErr(ginctx, &webbase.Error{
			StatusCode: http.StatusInternalServerError,
			Log:        fmt.Sprintf("failed to get org, err: %v", err),
		})
		return
	}
	err = w.manager.SetUserOrg(ginctx, email, org.Name)
	if err == nil {
		ginctx.Redirect(http.StatusMovedPermanently, "/org/manage")
		return
	}
	if errors.Is(err, auth.ErrNotFound) {
		w.ServeErr(ginctx, &webbase.Error{
			StatusCode: http.StatusBadRequest,
			Log:        fmt.Sprintf("failed to set user org; err: %v", err),
		})
		return
	}
	w.ServeErr(ginctx, &webbase.Error{
		StatusCode: http.StatusInternalServerError,
		Log:        fmt.Sprintf("failed to set user org; err: %v", err),
	})
}

// OrgRegister sets org.
func (w *Web) OrgRegister() gin.HandlerFunc {
	return w.Handler(
		"org.html.tmpl",
		func(ginctx *gin.Context) (interface{}, *webbase.Error) {
			return PageData{
				Data:           webbase.NewData("Golinks - Organization Creation", ginctx),
				FormInputName:  formInputName,
				FormInputEmail: formInputEmail,
				FormBtnAction:  formBtnActionCreate,
			}, nil
		},
	)
}

// HandleOrgRegisterForm handle request to create org.
func (w *Web) HandleOrgRegisterForm(ginctx *gin.Context) {
	user, err := ctx.GetUser(ginctx)
	if err != nil {
		w.ServeErr(ginctx, &webbase.Error{
			StatusCode: http.StatusInternalServerError,
			Log:        fmt.Sprintf("failed to get user, err: %v", err),
		})
		return
	}
	name := ginctx.PostForm(formInputName)
	org := auth.Organization{
		Name:       name,
		AdminEmail: user.Email,
	}
	err = w.manager.RegisterOrg(ginctx.Request.Context(), org)
	if err == nil {
		ginctx.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	if errors.Is(err, auth.ErrOrgExists) ||
		errors.Is(err, auth.ErrBadParams) {
		w.ServeErr(ginctx, &webbase.Error{
			StatusCode: http.StatusBadRequest,
			Log:        fmt.Sprintf("failed to register org. err: %v", err),
		})
		return
	}
	w.ServeErr(ginctx, &webbase.Error{
		StatusCode: http.StatusInternalServerError,
		Log:        fmt.Sprintf("failed to register org, err: %v", err),
	})
}
