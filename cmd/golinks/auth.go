package main

import (
	"strings"
	"time"

	"github.com/popodidi/log"

	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/auth/kv"
	"github.com/haostudio/golinks/internal/auth/traced"
	"github.com/haostudio/golinks/internal/encoding"
)

// AuthManagerConfig defines the auth provider config.
type AuthManagerConfig struct {
	TokenExpieration int    `conf:"default:30;usage:token expiration time in day"`
	TokenSecret      string `conf:"default:_golinks_jwt_token_secret_"`

	NoAuth struct {
		Enabled    bool   `conf:"default:false"`
		DefaultOrg string `conf:"default:_no_org_"`
	}

	Type string `conf:"default:kv"`
	Kv   StoreConfig
}

func newAuthManager(logger log.Logger,
	conf AuthManagerConfig, enc encoding.Binary, traceEnabled bool) (
	manager *auth.Manager, closeFunc func() error) {
	var provider auth.Provider
	if conf.NoAuth.Enabled {
		closeFunc = func() error { return nil }
		return
	}
	switch strings.ToLower(conf.Type) {
	case "kv":
		provider, closeFunc = newKvAuthProvider(logger, conf.Kv, enc, traceEnabled)
	default:
		logger.Critical("unknown auth provider type: %s", conf.Type)
	}
	if traceEnabled {
		provider = traced.New(provider)
	}
	manager = auth.New(auth.Config{
		Provider:         provider,
		TokenExpieration: time.Duration(conf.TokenExpieration) * 24 * time.Hour,
		TokenSecret:      []byte(conf.TokenSecret),
	})
	return
}

func newKvAuthProvider(logger log.Logger,
	conf StoreConfig, enc encoding.Binary, traceEnabled bool) (
	auth.Provider, func() error) {
	kvStore, closeFunc := newStore(logger, conf, traceEnabled)
	provider := kv.New(kvStore.In(authNamespace), enc)
	return provider, closeFunc
}
