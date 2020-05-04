package main

import (
	"fmt"
	"strings"

	"github.com/popodidi/log"

	"github.com/haostudio/golinks/internal/kv"
	"github.com/haostudio/golinks/internal/kv/bolt"
	"github.com/haostudio/golinks/internal/kv/cached"
	"github.com/haostudio/golinks/internal/kv/leveldb"
	"github.com/haostudio/golinks/internal/kv/memory"
	"github.com/haostudio/golinks/internal/kv/traced"
)

// StoreConfig defines a kv store config.
type StoreConfig struct {
	Type string `conf:"default:bolt"`

	// XXX: Use LRU as cache at the moment. In general, the cache store supports
	// any kind of store as cache. However, to reduce the config complexity at
	// the moment, we use only LRU as the cache store.
	LRUCache bool `conf:"default:true"`

	// Type-specific options
	Memory  MemStore
	Bolt    BoltStore
	LevelDB LevelDBStore
}

func newStore(logger log.Logger, conf StoreConfig, traceEnabled bool) (
	store kv.Store, closeFunc func() error) {
	// init canonical store
	closeCanonical := func() error { return nil }
	switch strings.ToLower(conf.Type) {
	case "leveldb":
		store = newLevelDBStore(logger, conf.LevelDB)
	case "bolt":
		store = newBoltStore(logger, conf.Bolt)
	case "memory":
		store = newMemStore(logger, conf.Memory, false)
	default:
		logger.Critical("unknown link storage type: %s", conf.Type)
	}

	closeCanonical = store.Close

	// init traced
	if traceEnabled {
		store = traced.New(store)
	}

	if !conf.LRUCache {
		closeFunc = closeCanonical
		return
	}

	closeCache := func() error { return nil }
	if conf.LRUCache {
		memConf := MemStore{
			Engine: "lru",
		}
		memConf.LRU.Cap = 1 << 10
		cacheStore := newMemStore(logger, memConf, true)
		closeCache = cacheStore.Close
		store = cached.New(store, cacheStore.In(cacheNamespace))
	}

	closeFunc = func() error {
		canonicalErr := closeCanonical()
		cacheErr := closeCache()
		if canonicalErr == nil &&
			cacheErr == nil {
			return nil
		}
		return fmt.Errorf("canonical:%v|cache:%v", canonicalErr, cacheErr)
	}
	return
}

// MemStore defines config for memory store.
type MemStore struct {
	Engine string `conf:"default:map"`
	LRU    struct {
		Cap int `conf:"default:1024"`
	}
}

func newMemStore(logger log.Logger, conf MemStore, asCache bool) kv.Store {
	if !asCache {
		logger.Warn("data in memory store is not persistent")
	}
	switch conf.Engine {
	case "map":
		return memory.New()
	case "lru":
		return memory.NewLRU(conf.LRU.Cap)
	default:
		logger.Critical("unsupported memory store engint: %s", conf.Engine)
		return nil
	}
}

// BoltStore defines config for bolt store.
type BoltStore struct {
	Dir  string `conf:"default:datadir"`    // directory for file system store
	Name string `conf:"default:golinks.db"` // db name for file system store
}

func newBoltStore(logger log.Logger, conf BoltStore) kv.Store {
	store, err := bolt.New(conf.Dir, conf.Name, rootNamespace)
	if err != nil {
		logger.Critical("failed to create bolt store. err: %v", err)
	}
	return store
}

// LevelDBStore defines config for leveldb store.
type LevelDBStore struct {
	Dir  string `conf:"default:datadir"`    // directory for file system store
	Name string `conf:"default:golinks.db"` // db name for file system store
}

func newLevelDBStore(logger log.Logger, conf LevelDBStore) kv.Store {
	logger.Warn("leveldb store may run with concurrency issues")
	logger.Warn("use with your own risk or use bolt store instead")
	store, err := leveldb.New(conf.Dir, conf.Name, nil)
	if err != nil {
		logger.Critical("failed to create bolt store. err: %v", err)
	}
	return store
}
