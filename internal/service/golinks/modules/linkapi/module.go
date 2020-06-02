package linkapi

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/api/middlewares"
	"github.com/haostudio/golinks/internal/link"
	"github.com/haostudio/golinks/internal/service/golinks/ctx"
)

// New returns a new link api module.
func New(store link.Store) *Links {
	return &Links{
		store: store,
	}
}

// Links defines the link module struct.
type Links struct {
	store link.Store
}

// PathParamLinkKey returns the link_key path parameter.
func (l *Links) PathParamLinkKey() string {
	return "link_key"
}

// GetLinks returns all links.
func (l *Links) GetLinks(ginctx *gin.Context) {
	logger := middlewares.GetLogger(ginctx)

	// get links
	org, err := ctx.GetOrg(ginctx)
	if err != nil {
		logger.Error("failed to get org. err: %v", err)
		ginctx.Status(http.StatusInternalServerError)
		return
	}
	links, err := l.store.GetLinks(ginctx.Request.Context(), org.Name)
	if err != nil {
		logger.Error("failed to get links from store. err: %v", err)
		ginctx.Status(http.StatusInternalServerError)
		return
	}

	// sorted keys
	keys := make([]string, len(links))
	var i int
	for k := range links {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// construct response
	res := make(map[string]string)
	for _, key := range keys {
		ln := links[key]
		desc, err := ln.Description()
		if err != nil {
			logger.Debug(
				"failed to get description of link with key \"%s\". err: %v", key, err)
			continue
		}
		res[key] = desc
	}

	ginctx.JSON(http.StatusOK, res)
}

// UpdateLink updates the link.
func (l *Links) UpdateLink(ginctx *gin.Context) {
	logger := middlewares.GetLogger(ginctx)
	key := ginctx.Param(l.PathParamLinkKey())
	if len(key) == 0 {
		logger.Error("empty key")
		ginctx.String(http.StatusBadRequest, "empty key")
		return
	}

	// read link from request
	var req struct {
		Version int    `json:"version"`
		Payload string `json:"payload"`
	}
	err := ginctx.BindJSON(&req)
	if err != nil {
		logger.Error("failed to bind json. err: %v", err)
		return
	}
	link, err := link.New(req.Version, req.Payload)
	if err != nil {
		logger.Error("failed to bind json. err: %v", err)
		ginctx.Status(http.StatusBadRequest)
		return
	}

	// update to store
	org, err := ctx.GetOrg(ginctx)
	if err != nil {
		logger.Error("failed to get org. err: %v", err)
		ginctx.Status(http.StatusInternalServerError)
		return
	}
	err = l.store.UpdateLink(ginctx.Request.Context(), org.Name, key, *link)
	if err != nil {
		logger.Error(
			"failed to update \"%s\" to %s in store. err: %v", key, *link, err)
		ginctx.Status(http.StatusInternalServerError)
		return
	}
	ginctx.Status(http.StatusOK)
}

// DeleteLink deletes the link.
func (l *Links) DeleteLink(ginctx *gin.Context) {
	logger := middlewares.GetLogger(ginctx)
	key := ginctx.Param(l.PathParamLinkKey())
	if len(key) == 0 {
		logger.Error("empty key")
		ginctx.String(http.StatusBadRequest, "empty key")
		return
	}

	org, err := ctx.GetOrg(ginctx)
	if err != nil {
		logger.Error("failed to get org. err: %v", err)
		ginctx.Status(http.StatusInternalServerError)
		return
	}
	err = l.store.DeleteLink(ginctx.Request.Context(), org.Name, key)
	if err != nil {
		logger.Error("failed to delete \"%s\" from store. err: %v", key, err)
		ginctx.Status(http.StatusInternalServerError)
		return
	}
	ginctx.Status(http.StatusOK)
}
