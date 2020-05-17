package webbase

import (
	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/version"
)

// Data defines the base data.
type Data struct {
	Title      string
	LoggedIn   bool
	BinVersion string
	Ctx        struct {
		Org, User string
		LoggedIn  bool
	}
}

// NewData returns a database.
func NewData(title string, ctx *gin.Context) Data {
	data := Data{
		Title:      title,
		BinVersion: version.Version(),
	}
	org, err := middlewares.GetOrg(ctx)
	if err != nil {
		return data
	}
	data.Ctx.Org = org.Name
	user, err := middlewares.GetUser(ctx)
	if err != nil {
		return data
	}
	data.Ctx.User = user.Email
	data.Ctx.LoggedIn = true
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
