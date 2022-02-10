package tracing

import (
	"io"

	"trellis.tech/trellis.v1/pkg/trellis"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

func NewJeagerTracer(name string, cfg *trellis.TracingConfig) (opentracing.Tracer, io.Closer, error) {

	metricsFactory := prometheus.New()
	return config.Configuration{
		Sampler:     &config.SamplerConfig{},
		ServiceName: name,
		Disabled:    cfg.Enable,
	}.NewTracer(
		config.Metrics(metricsFactory),
	)
}
