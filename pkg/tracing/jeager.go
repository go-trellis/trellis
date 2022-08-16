/*
Copyright © 2022 Henry Huang <hhh@rutcode.com>
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package tracing

import (
	"io"

	"trellis.tech/trellis.v1/pkg/trellis"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	jeager "github.com/uber/jaeger-client-go"
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

func GetTraceID(span opentracing.Span) string {
	if span == nil {
		return uuid.New().String()
	}
	switch t := span.(type) {
	case *jeager.Span:
		return t.SpanContext().TraceID().String()
	default:
		return uuid.New().String()
	}
}

func GetSpanID(span opentracing.Span) string {
	if span == nil {
		return uuid.New().String()
	}
	switch t := span.(type) {
	case *jeager.Span:
		return t.SpanContext().SpanID().String()
	default:
		return uuid.New().String()
	}
}

func GetSpanContextID(span opentracing.Span) string {
	if span == nil {
		return uuid.New().String()
	}
	switch t := span.(type) {
	case *jeager.Span:
		return t.SpanContext().String()
	default:
		return uuid.New().String()
	}
}
