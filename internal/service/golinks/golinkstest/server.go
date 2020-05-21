package golinkstest

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/internal/auth"
	authkv "github.com/haostudio/golinks/internal/auth/kv"
	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/kv/memory"
	lnkv "github.com/haostudio/golinks/internal/link/kv"
	"github.com/haostudio/golinks/internal/service"
	"github.com/haostudio/golinks/internal/service/golinks"
)

// NewTestServer returns a server with memory store and cache.
func NewTestServer() service.Service {
	store := memory.New()
	enc := gob.New()

	lnStore := lnkv.New(store.In("link"), enc)
	authProvider := authkv.New(store.In("auth"), enc)

	conf := golinks.Config{
		Gin:       gin.Default(),
		Address:   "0.0.0.0:8000",
		Traced:    false,
		LinkStore: lnStore,
	}
	conf.Auth.Enabled = true
	conf.Auth.DefaultOrg = ""
	conf.Auth.Manager = auth.New(auth.Config{
		Provider:         authProvider,
		TokenExpieration: 1 * 24 * time.Hour,
		TokenSecret:      []byte("token_secret"),
	})

	gin.Default()

	return golinks.New(conf)
}
