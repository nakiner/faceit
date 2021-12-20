package metrics

import (
	"context"
	"sync"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var once = sync.Once{}

type Metrics struct {
	Counter    metrics.Counter
	Histogramm metrics.Histogram
}

type metricKey struct{}

func (m *Metrics) Get() (metrics.Counter, metrics.Histogram) {
	once.Do(func() {
		fieldKeys := []string{"handler", "code", "service"}
		m.Counter = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "faceit",
			Name:      "requests_total",
			Help:      "Number of requests received.",
		}, fieldKeys)
		m.Histogramm = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "faceit",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys)
	})

	return m.Counter, m.Histogramm
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, metricKey{}, new(Metrics))
}

func FromContext(ctx context.Context) (metrics.Counter, metrics.Histogram) {
	if metric, ok := ctx.Value(metricKey{}).(*Metrics); ok {
		return metric.Get()
	}
	return nil, nil
}
