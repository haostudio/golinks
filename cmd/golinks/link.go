package main // nolint: dupl

import (
	"strings"

	"github.com/popodidi/log"

	"github.com/haostudio/golinks/internal/encoding"
	"github.com/haostudio/golinks/internal/link"
	"github.com/haostudio/golinks/internal/link/kv"
	"github.com/haostudio/golinks/internal/link/traced"
)

// LinkStoreConfig defines the link store config.
type LinkStoreConfig struct {
	Type string `conf:"default:kv"`
	Kv   StoreConfig
}

func newLinkStore(logger log.Logger,
	conf LinkStoreConfig, enc encoding.Binary, traceEnabled bool) (
	store link.Store, closeFunc func() error) {
	switch strings.ToLower(conf.Type) {
	case "kv":
		store, closeFunc = newKvLinkStore(logger, conf.Kv, enc, traceEnabled)
	default:
		logger.Critical("unknown link store type: %s", conf.Type)
	}
	if traceEnabled {
		store = traced.New(store)
	}
	return
}

func newKvLinkStore(logger log.Logger,
	conf StoreConfig, enc encoding.Binary, traceEnabled bool) (
	link.Store, func() error) {
	linkKv, closeFunc := newStore(logger, conf, traceEnabled)
	store := kv.New(linkKv.In(linkNamespace), enc)
	return store, closeFunc
}
