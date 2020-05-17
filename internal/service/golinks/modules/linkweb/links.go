package linkweb

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/link"
	"github.com/haostudio/golinks/internal/service/golinks/modules/webbase"
)

const (
	formInputVersion = "version"
	formInputPayload = "payload"
	formInputAction  = "action"
	formSaveValue    = "Save"
	formDeleteValue  = "Delete"
)

// Config defines the web config.
type Config struct {
	Store  link.Store
	Traced bool
}

// Web defines the web handler module.
type Web struct {
	webbase.Base
	store link.Store
}

// New returns a new web handler module.
func New(conf Config) *Web {
	return &Web{
		Base:  webbase.NewBase(conf.Traced),
		store: conf.Store,
	}
}

// PathParamLinkKey returns the link_key path parameter.
func (w *Web) PathParamLinkKey() string {
	return "link_key"
}

// Links returns the page for the list of links. (./web/links.html)
func (w *Web) Links() gin.HandlerFunc {
	return w.Handler(
		"links.html.tmpl",
		func(ctx *gin.Context) (interface{}, *webbase.Error) {
			// get links
			org, err := middlewares.GetOrg(ctx)
			if err != nil {
				return nil, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
					Log:        fmt.Sprintf("failed to get org. err: %v", err),
				}
			}
			links, err := w.store.GetLinks(ctx.Request.Context(), org.Name)
			if errors.Is(err, link.ErrNotFound) {
				links = make(map[string]link.Link)
			} else if err != nil {
				return nil, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
					Log: fmt.Sprintf(
						"failed to get links from store. err: %v", err,
					),
				}
			}

			// sorted keys
			keys := make([]string, len(links))
			var i int
			for k := range links {
				keys[i] = k
				i++
			}
			sort.Strings(keys)

			// construct data
			logger := middlewares.GetLogger(ctx)
			pageData := NewAllPageData(ctx)
			for _, key := range keys {
				ln := links[key]
				lnData, err := NewLink(key, ln)
				if err != nil {
					logger.Error("failed to get links data of \"%s\". err: %v", key, err)
					continue
				}
				pageData.Links = append(pageData.Links, lnData)
			}
			return pageData, nil
		},
	)
}

// EditLink returns the edit page of a link (./web/edit.yaml)
func (w *Web) EditLink() gin.HandlerFunc {
	return w.Handler(
		"edit.html.tmpl",
		func(ctx *gin.Context) (interface{}, *webbase.Error) {
			key := ctx.Param(w.PathParamLinkKey())

			pageData := NewEditPageData(ctx)
			pageData.FormInputVersion = formInputVersion
			pageData.FormInputPayload = formInputPayload
			pageData.FormInputAction = formInputAction
			pageData.FormSaveValue = formSaveValue
			pageData.FormDeleteValue = formDeleteValue
			pageData.Link.Key = key

			org, err := middlewares.GetOrg(ctx)
			if err != nil {
				return nil, &webbase.Error{
					StatusCode: http.StatusInternalServerError,
					Log:        fmt.Sprintf("failed to get org. err: %v", err),
				}
			}
			for len(key) != 0 {
				// Get link from store
				ln, err := w.store.GetLink(ctx.Request.Context(), org.Name, key)
				if errors.Is(err, link.ErrNotFound) {
					break
				}
				if err != nil {
					return nil, &webbase.Error{
						StatusCode: http.StatusInternalServerError,
						Log: fmt.Sprintf(
							"failed to get link from store. err: %v", err,
						),
					}
				}
				pageData.Link, err = NewLink(key, ln)
				if err != nil {
					return nil, &webbase.Error{
						StatusCode: http.StatusInternalServerError,
						Log: fmt.Sprintf(
							"failed to get links data of \"%s\". err: %v", key, err,
						),
					}
				}
				break
			}
			return pageData, nil
		},
	)
}

// HandleEditLinktForm handles the edit.html form submission.
func (w *Web) HandleEditLinktForm(ctx *gin.Context) {
	key := ctx.Param(w.PathParamLinkKey())
	if len(key) == 0 {
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusInternalServerError,
			Messages:   []string{"Empty key"},
			Log:        "Empty key",
		})
		return
	}

	action := ctx.PostForm("action")
	version := ctx.PostForm("version")
	payload := ctx.PostForm("payload")

	org, err := middlewares.GetOrg(ctx)
	if err != nil {
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusInternalServerError,
			Log:        fmt.Sprintf("failed to get org. err: %v", err),
		})
		return
	}
	switch action {
	case formSaveValue:
		v, err := strconv.Atoi(version)
		if err != nil {
			w.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log: fmt.Sprintf(
					"failed to parse version %s. err: %v", version, err,
				),
			})
			return
		}
		ln, err := link.New(v, payload)
		if err != nil {
			w.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log: fmt.Sprintf(
					"failed to create link from request form. err: %v", err,
				),
			})
			return
		}
		// update to store
		err = w.store.UpdateLink(ctx.Request.Context(), org.Name, key, *ln)
		if err != nil {
			w.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log: fmt.Sprintf(
					"failed to update \"%s\" to %v in store. err: %v", key, *ln, err),
			})
			return
		}
		ctx.Redirect(http.StatusMovedPermanently, "/links")
		return
	case formDeleteValue:
		err := w.store.DeleteLink(ctx.Request.Context(), org.Name, key)
		if err != nil {
			w.ServeErr(ctx, &webbase.Error{
				StatusCode: http.StatusInternalServerError,
				Log: fmt.Sprintf(
					"failed to delete \"%s\" from store. err: %v", key, err),
			})
			return
		}
		ctx.Redirect(http.StatusMovedPermanently, "/links")
		return
	default:
		w.ServeErr(ctx, &webbase.Error{
			StatusCode: http.StatusBadRequest,
			Messages:   []string{"Invalid action"},
			Log:        fmt.Sprintf("invalid action %s", action),
		})
		return
	}
}
