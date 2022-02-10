package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"trellis.tech/trellis.v1/pkg/trellis"
)

func InitTracer(name string, cfg *trellis.TracingConfig) (io.Closer, error) {
	if cfg == nil {
		return nil, nil
	}

	tracer, closer, err := NewJeagerTracer(name, cfg)
	if err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
