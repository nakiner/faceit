package sentry

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/nakiner/faceit/configs"
	"github.com/nakiner/faceit/pkg/health"
)

// NewSentry initializes new global sentry Client.
func NewSentry(cfg *configs.Config) error {
	var debug bool
	if cfg.Logger.Level == "debug" {
		debug = true
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Sentry.Dsn,
		Debug:            debug,
		AttachStacktrace: true,
		Environment:      cfg.Sentry.Environment,
		Release:          health.Commit,
	}); err != nil {
		return err
	}
	var data map[string]interface{}
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: "config",
		Data:     data,
		Message:  "init config",
	})
	return nil
}

func Middleware(h http.Handler) http.Handler {
	return sentryhttp.New(sentryhttp.Options{Repanic: true}).Handle(h)
}
