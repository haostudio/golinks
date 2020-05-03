package main

import (
	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers/iowriter"
	"github.com/popodidi/log/handlers/multi"
)

// LogConfig defines the log config struct.
type LogConfig struct {
	Level  int `conf:"default:6"`
	Stdout struct {
		Enabled   bool `conf:"default:true"`
		WithColor bool `conf:"default:true"`
	}
}

func configLogger(config LogConfig) (closeFunc func() error) {
	handler := multi.New()
	if config.Stdout.Enabled {
		handler = multi.New(handler, iowriter.Stdout(config.Stdout.WithColor))
	}
	log.Set(log.Config{
		Threshold: log.Level(config.Level),
		Handler:   handler,
	})
	return handler.Close
}
