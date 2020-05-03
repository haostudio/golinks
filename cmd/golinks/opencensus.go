package main

import (
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	"github.com/popodidi/log"
)

// MetricsConfig defines the metrics config struct.
type MetricsConfig struct {
	SampleRate float64 `conf:"default:1"`
	Jaeger     struct {
		Enabled           bool   `conf:"default:false"`
		AgentEndpoint     string `conf:"default:localhost:6831"`
		CollectorEndpoint string `conf:"default:http://localhost:14268/api/traces"`
	}
}

// Enabled returns if any metrics is enabled.
func (c MetricsConfig) Enabled() bool {
	return c.Jaeger.Enabled
}

func configMetrics(logger log.Logger, conf MetricsConfig) (closeFunc func()) {
	closeFunc = func() {}
	if !conf.Enabled() {
		return
	}

	view.SetReportingPeriod(time.Second)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(conf.SampleRate),
	})

	if conf.Jaeger.Enabled {
		je, err := jaeger.NewExporter(jaeger.Options{
			AgentEndpoint:     conf.Jaeger.AgentEndpoint,
			CollectorEndpoint: conf.Jaeger.CollectorEndpoint,
			ServiceName:       "golinks",
		})
		if err != nil {
			logger.Critical("failed to create the jaeger exporter: %v", err)
		}
		trace.RegisterExporter(je)
	}
	return
}
