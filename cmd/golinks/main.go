package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/popodidi/conf"
	"github.com/popodidi/conf/source/env"
	"github.com/popodidi/conf/source/yaml"
	"github.com/popodidi/log"

	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/service"
	"github.com/haostudio/golinks/internal/service/golinks"
	"github.com/haostudio/golinks/internal/version"
)

const (
	rootNamespace  = "github.com/haostudio/golinks"
	linkNamespace  = "_link"
	authNamespace  = "_auth"
	cacheNamespace = "_cache"
)

// Config defines golinks server config.
type Config struct {
	Log          LogConfig
	Port         int `conf:"default:8000"`
	Metrics      MetricsConfig
	LinkStore    LinkStoreConfig
	AuthProvider AuthProviderConfig
	HTTP         struct {
		Golinks struct {
			Enabled bool `conf:"default:true"`
			Wiki    bool `conf:"default:false"`
		}
	}
}

func main() {
	// Load server config
	var config Config
	cfg := conf.New(&config)
	err := cfg.Load(env.New(), yaml.New("golinks_config.yaml"))
	if err != nil {
		panic(err)
	}

	// Configure logger
	closeFunc := configLogger(config.Log)
	defer func() {
		err := closeFunc()
		if err != nil {
			panic(err)
		}
	}()

	logger := log.New("golinks")

	// Print version and config
	logGolinksInfo(logger, cfg)

	// Catch panic
	defer func(logger log.Logger) {
		if recovered := recover(); recovered != nil {
			logger.Error("%v", recovered)
			logger.Error("%s", debug.Stack())
			panic(recovered)
		}
	}(logger)

	// Configure metrics
	closeMetrics := configMetrics(logger, config.Metrics)
	defer closeMetrics()

	// Golinks uses gob encoding
	enc := gob.New()

	// links store
	linkStore, linkStoreClose := newLinkStore(
		logger, config.LinkStore, enc, config.Metrics.Enabled(),
	)
	defer func() {
		err := linkStoreClose()
		if err != nil {
			logger.Warn("failed to close link store. %v", err)
		}
	}()

	// auth provider
	authProvider, authProviderClose := newAuthProvider(
		logger, config.AuthProvider, enc, config.Metrics.Enabled(),
	)
	defer func() {
		err := authProviderClose()
		if err != nil {
			logger.Warn("failed to close auth provider. %v", err)
		}
	}()

	// Setup service mux
	mux := service.NewMux(logger)
	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)

	// Setup default HTTP server
	if config.HTTP.Golinks.Enabled {
		golinksConfig := golinks.Config{
			Gin:       gin.New(),
			Address:   addr,
			Traced:    config.Metrics.Enabled(),
			Wiki:      config.HTTP.Golinks.Wiki,
			LinkStore: linkStore,
		}
		golinksConfig.Auth.Enabled = !config.AuthProvider.NoAuth.Enabled
		golinksConfig.Auth.DefaultOrg = config.AuthProvider.NoAuth.DefaultOrg
		golinksConfig.Auth.Provider = authProvider
		mux.Append(golinks.New(golinksConfig))
	}

	logger.Info("server listening to port [%d]", config.Port)
	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Critical("failed to listen TCP. %v", err)
	}
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)
		sig := <-quitChan

		logger.Warn("signal [%s] received. shutting down TCP server...", sig)
		if err := tcpListener.Close(); err != nil {
			logger.Error("tcp listener stopped with error %v", err)
		}
	}()
	err = mux.Serve(tcpListener)
	if err != nil {
		logger.Warn("tcp mux stopped: %v", err)
	}
	logger.Info("process terminated")
}

func logGolinksInfo(logger log.Logger, cfg *conf.Config) {
	// Print version info
	ver := fmt.Sprintf("HAOSTUDIO/GOLINKS %s", version.Version())
	bar := strings.Repeat("=", len(ver))
	logger.Info(bar)
	logger.Info(ver)
	logger.Info(bar)

	// Print config
	configMap, err := cfg.Map()
	if err != nil {
		logger.Critical("failed to get config map. err: %v", err)
	}
	configMap.Iter(func(key string, val interface{}, path ...string) (next bool) {
		if len(path) == 0 {
			logger.Debug("%s: %v", key, val)
		} else {
			logger.Debug("%s.%s: %v", strings.Join(path, "."), key, val)
		}
		return true
	})
	logger.Debug(bar)
}
