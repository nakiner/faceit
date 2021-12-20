package tracing

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nakiner/faceit/tools/logging"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type tracerKey struct{}

func WithContext(ctx context.Context, tracer opentracing.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

func FromContext(ctx context.Context) opentracing.Tracer {
	if tracer, ok := ctx.Value(tracerKey{}).(opentracing.Tracer); ok {
		return tracer
	}
	return opentracing.GlobalTracer()
}

func NewJaegerTracer(ctx context.Context, addr, name string) (opentracing.Tracer, io.Closer, error) {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "component", "tracer")
	tracer, closer, err := jaegercfg.Configuration{
		ServiceName: name,
		Sampler: &jaegercfg.SamplerConfig{
			SamplingServerURL: addr,
			Type:              jaeger.SamplerTypeConst,
			Param:             1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort:  addr,
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}.NewTracer(
		jaegercfg.Logger(&jaegerLoggerAdapter{logger}),
	)
	if err != nil {
		return nil, ioutil.NopCloser(nil), err
	}
	return tracer, closer, nil
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (la *jaegerLoggerAdapter) Error(msg string) {
	level.Error(la.logger).Log("err", msg)
}

func (la *jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	level.Debug(la.logger).Log("msg", fmt.Sprintf(msg, args...))
}
