package limiting

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nakiner/faceit/tools/logging"
	"golang.org/x/time/rate"
)

type Limiter interface {
	Middleware(next http.Handler) http.Handler
}

func NewLimiter(ctx context.Context, limit float64) Limiter {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "component", "limiter")
	return &limiter{
		limiter: rate.NewLimiter(rate.Limit(limit), 1),
		logger:  logger,
	}
}

type limiter struct {
	limiter *rate.Limiter
	logger  log.Logger
}

func (l *limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l.limiter.Allow() == false {
			level.Debug(l.logger).Log(
				"code", http.StatusTooManyRequests,
				"msg", http.StatusText(http.StatusTooManyRequests),
			)
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
