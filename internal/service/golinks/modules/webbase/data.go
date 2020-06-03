package webbase

import (
	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/service/golinks/ctx"
	"github.com/haostudio/golinks/internal/version"
)

// Data defines the base data.
type Data struct {
	Title      string
	BinVersion string
	Ctx        struct {
		Org, User   string
		LoggedIn    bool
		AuthEnabled bool
	}
}

// NewData returns a database.
func NewData(title string, ginctx *gin.Context) Data {
	data := Data{
		Title:      title,
		BinVersion: version.Version(),
	}
	org, err := ctx.GetOrg(ginctx)
	if err == nil {
		data.Ctx.Org = org.Name
	}
	user, err := ctx.GetUser(ginctx)
	if err == nil {
		data.Ctx.User = user.Email
		data.Ctx.LoggedIn = true
	}
	data.Ctx.AuthEnabled = ctx.IsAuthEnabled(ginctx)
	return data
}

// ErrorPage defines the data for error.html template.
type ErrorPage struct {
	Data
	StatusCode  int
	ErrTitle    string
	Description string
}

// NewErrorPage returns error page data.
func NewErrorPage(ctx *gin.Context) ErrorPage {
	return ErrorPage{
		Data: NewData("Golinks - Error", ctx),
	}
}
