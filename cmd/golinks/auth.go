package main // nolint: dupl

import (
	"strings"

	"github.com/popodidi/log"

	"github.com/haostudio/golinks/internal/auth"
	"github.com/haostudio/golinks/internal/auth/kv"
	"github.com/haostudio/golinks/internal/auth/traced"
	"github.com/haostudio/golinks/internal/encoding"
)

// AuthProviderConfig defines the auth provider config.
type AuthProviderConfig struct {
	NoAuth struct {
		Enabled    bool   `conf:"default:false"`
		DefaultOrg string `conf:"default:_no_org_"`
	}
	Type string `conf:"default:kv"`
	Kv   StoreConfig
}

func newAuthProvider(logger log.Logger,
	conf AuthProviderConfig, enc encoding.Binary, traceEnabled bool) (
	provider auth.Provider, closeFunc func() error) {
	if conf.NoAuth.Enabled {
		provider = nil
		closeFunc = func() error { return nil }
		return
	}
	switch strings.ToLower(conf.Type) {
	case "kv":
		provider, closeFunc = newKvAuthProvider(logger, conf.Kv, enc, traceEnabled)
	default:
		logger.Critical("unknown link store type: %s", conf.Type)
	}
	if traceEnabled {
		provider = traced.New(provider)
	}
	return
}
func newKvAuthProvider(logger log.Logger,
	conf StoreConfig, enc encoding.Binary, traceEnabled bool) (
	auth.Provider, func() error) {
	kvStore, closeFunc := newStore(logger, conf, traceEnabled)
	provider := kv.New(kvStore.In(authNamespace), enc)
	return provider, closeFunc
}
