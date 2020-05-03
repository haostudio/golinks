package servertest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haostudio/golinks/cmd/golinks/server"
	authkv "github.com/haostudio/golinks/internal/auth/kv"
	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/kv/memory"
	lnkv "github.com/haostudio/golinks/internal/link/kv"
)

// NewTestServer returns a server with memory store and cache.
func NewTestServer() http.Handler {
	store := memory.New()
	enc := gob.New()

	lnStore := lnkv.New(store.In("link"), enc)
	authProvider := authkv.New(store.In("auth"), enc)

	conf := server.Config{
		Gin:          gin.Default(),
		AuthProvider: authProvider,
		LinkStore:    lnStore,
	}
	gin.Default()

	return server.New(conf)
}
