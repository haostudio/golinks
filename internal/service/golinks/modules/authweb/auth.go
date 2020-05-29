package authweb

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

const (
	formInputEmail        = "email"
	formInputPassword     = "password"
	formInputName         = "name"
	formInputCallback     = "callback"
	formBtnAction         = "save"
	formInputAction       = "action"
	formRegisterBtnAction = "register"
	formLoginBtnAction    = "login"
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
		func(ctx *gin.Context) (interface{}, *webbase.Error) {
			callback := ctx.Query("callback")
			if callback == "" {
				callback = "/"
			}
			return LoginData{
				Data:                  webbase.NewData("Golinks - Login", ctx),
				FormInputEmail:        formInputEmail,
				FormInputPassword:     formInputPassword,
				FormInputCallback:     formInputCallback,
				FormInputAction:       formInputAction,
				FormLoginBtnAction:    formLoginBtnAction,
				FormRegisterBtnAction: formRegisterBtnAction,
				Callback:              callback,
			}, nil
		},
	)
}

// HandleLoginForm handle request to create org user.
func (w *Web) HandleLoginForm(ctx *gin.Context) {
	email := ctx.PostForm(formInputEmail)
	password := ctx.PostForm(formInputPassword)
	action := ctx.PostForm(formInputAction)
	callback := ctx.PostForm(formInputCallback)
	switch action {
	case formLoginBtnAction:
		token, err := w.manager.Login(ctx.Request.Context(), email, password)
		if err != nil {
			w.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusUnauthorized,
				Messages:   []string{"Invalid email/password"},
				Log:        fmt.Sprintf("login failed. %v", err),
			})
			return
		}
		// authorized
		middlewares.SetToken(
			ctx, token.JWT, int(w.manager.TokenExpieration.Seconds()))
		ctx.Redirect(http.StatusMovedPermanently, callback)
		return
	case formRegisterBtnAction:
		fallthrough
	default:
		w.ServeErr(ctx, &webbase.Error{
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
		func(ctx *gin.Context) (interface{}, *webbase.Error) {
			org, err := middlewares.GetOrg(ctx)
			if err != nil {
				return nil, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
					Log:        fmt.Sprintf("failed to get org, err: %v", err),
				}
			}

			user, err := middlewares.GetUser(ctx)
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

			users, err := w.manager.GetOrgUsers(ctx, org.Name)
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
					"Golinks - Organization Management", ctx),
				Users:             userSlice,
				Admin:             admin,
				FormInputEmail:    formInputEmail,
				FormInputPassword: formInputPassword,
				FormBtnAction:     formBtnAction,
			}, nil
		},
	)
}

// HandleSetOrgUserForm handle request to create org user.
func (w *Web) HandleSetOrgUserForm(ctx *gin.Context) {
	email := ctx.PostForm(formInputEmail)
	password := ctx.PostForm(formInputPassword)
	org, err := middlewares.GetOrg(ctx)
	if err != nil {
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusInternalServerError,
			Log:        fmt.Sprintf("failed to get org, err: %v", err),
		})
		return
	}
	user, err := auth.NewUser(email, password, org.Name)
	if err != nil {
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusInternalServerError,
			Log:        fmt.Sprintf("failed to init user, err: %v", err),
		})
		return
	}
	err = w.manager.RegisterUser(ctx, *user)
	if err == nil {
		ctx.Redirect(http.StatusMovedPermanently, "/org/manage")
		return
	}
	if errors.Is(err, auth.ErrUserExists) {
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusBadRequest,
			Log:        "user exists",
		})
		return
	}
	if errors.Is(err, auth.ErrNotFound) {
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusBadRequest,
			Log:        "org not found",
		})
		return
	}
	w.ServeErr(ctx, &webbase.Error{
		StatusCode: http.StatusInternalServerError,
		Log:        fmt.Sprintf("failed to set user, err: %v", err),
	})
}

// SetOrg sets org.
func (w *Web) SetOrg() gin.HandlerFunc {
	return w.Handler(
		"org.html.tmpl",
		func(ctx *gin.Context) (interface{}, *webbase.Error) {
			return PageData{
				Data:              webbase.NewData("Golinks - Organization Creation", ctx),
				FormInputName:     formInputName,
				FormInputEmail:    formInputEmail,
				FormInputPassword: formInputPassword,
				FormBtnAction:     formBtnAction,
			}, nil
		},
	)
}

// HandleSetOrgForm handle request to create org.
func (w *Web) HandleSetOrgForm(ctx *gin.Context) {
	name := ctx.PostForm(formInputName)
	email := ctx.PostForm(formInputEmail)
	password := ctx.PostForm(formInputPassword)
	org := auth.Organization{
		Name:       name,
		AdminEmail: email,
	}
	user, err := auth.NewUser(email, password, name)
	if err != nil {
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusInternalServerError,
			Log:        fmt.Sprintf("failed to init user, err: %v", err),
		})
		return
	}
	err = w.manager.RegisterOrgWithAdmin(ctx.Request.Context(), org, *user)
	if err == nil {
		// In order to redirect to user add page with authorization,
		// we hack redirect url with credential.
		ctx.Writer.Header().Set(
			"Location",
			fmt.Sprintf(
				"http://%s:%s@%s/org/manage",
				email,
				password,
				ctx.Request.Host,
			),
		)
		ctx.AbortWithStatus(http.StatusMovedPermanently)
		return
	}
	if errors.Is(err, auth.ErrOrgExists) {
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusBadRequest,
			Log:        "org already exists",
		})
		return
	}
	w.ServeErr(ctx, &webbase.Error{
		StatusCode: http.StatusInternalServerError,
		Log:        fmt.Sprintf("failed to register org, err: %v", err),
	})
}
