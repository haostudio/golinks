package linkweb

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/link"
	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

// Link defines a link data for template
type Link struct {
	Exists  bool
	Key     string
	Version int
	Format  string
}

// NewLink returns a new link data.
func NewLink(key string, ln link.Link) (data Link, err error) {
	data.Exists = true
	data.Key = key
	data.Version = ln.Version
	desc, err := ln.Description()
	if err != nil {
		err = fmt.Errorf(
			"failed to get description of link with key \"%s\". err: %w",
			key, err,
		)
		return
	}
	// XXX: Since link description is now "vx|...|http(s)://xxx", we simply trim
	// the leading "vx|...|" and return the right.
	_, desc = link.Pop(desc, "http")
	desc = "http" + desc
	data.Format = desc
	return
}

// AllPageData defines the data for links.html template.
type AllPageData struct {
	webbase.Data
	Links []Link
}

// NewAllPageData returns links page data.
func NewAllPageData(ctx *gin.Context) AllPageData {
	return AllPageData{
		Data: webbase.NewData("Golinks - All links", ctx),
	}
}

// EditPageData defines the data for edit.html template.
type EditPageData struct {
	webbase.Data

	FormInputVersion string
	FormInputPayload string
	FormInputAction  string
	FormSaveValue    string
	FormDeleteValue  string

	Link Link
}

// NewEditPageData returns edit page data.
func NewEditPageData(ctx *gin.Context) EditPageData {
	return EditPageData{
		Data: webbase.NewData("Golinks - Edit links", ctx),
	}
}
