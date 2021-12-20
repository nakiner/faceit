package main

import (
	"context"
	"fmt"
	"github.com/nakiner/faceit/internal/store/database"
	"github.com/nakiner/faceit/pkg/health"
	userQueue "github.com/nakiner/faceit/pkg/queue/user"
	natsCl "github.com/nakiner/faceit/pkg/store/nats"
	"github.com/nakiner/faceit/pkg/user"
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/nakiner/faceit/configs"
	"github.com/nakiner/faceit/internal/server"
	"github.com/nakiner/faceit/tools/logging"
	"github.com/nakiner/faceit/tools/metrics"
	"github.com/nakiner/faceit/tools/sentry"
	"github.com/nakiner/faceit/tools/tracing"

	userRepository "github.com/nakiner/faceit/internal/repository/user"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config
	cfg := configs.NewConfig()
	if err := cfg.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}
	// Print config
	if err := cfg.Print(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}

	logger, err := logging.NewLogger(cfg.Logger.Level, cfg.Logger.TimeFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %s", err)
		os.Exit(1)
	}
	ctx = logging.WithContext(ctx, logger)

	if cfg.Tracer.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(
			ctx,
			fmt.Sprintf("%s:%d", cfg.Tracer.Host, cfg.Tracer.Port),
			cfg.Tracer.Name,
		)
		if err != nil {
			level.Error(logger).Log("err", err, "msg", "failed to init tracer")
		}
		defer closer.Close()
		ctx = tracing.WithContext(ctx, tracer)
	}
	if cfg.Sentry.Enabled {
		if err := sentry.NewSentry(cfg); err != nil {
			level.Error(logger).Log("err", err, "msg", "failed to init sentry")
		}
	}

	if cfg.Metrics.Enabled {
		ctx = metrics.WithContext(ctx)
	}

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		level.Error(logger).Log("msg", "failed to init db", "err", err)
		os.Exit(1)
	}

	defer db.Close()

	nc, err := natsCl.NewClient(&cfg.Nats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init nats: %s", err)
		os.Exit(1)
	}
	defer nc.Close()

	ec, err := natsCl.NewEncodedClient(nc)
	if err != nil {
		level.Error(logger).Log("msg", "err init nats NewEncodedConn", "err", err)
		os.Exit(1)
	}
	defer ec.Close()

	userRepo := initUserRepository(ctx, db, cfg)
	userNatsPub, err := userQueue.NewPublisher(ec)
	if err != nil {
		level.Error(logger).Log("msg", "err init userQueue.Publisher", "err", err)
		os.Exit(1)
	}

	healthService := initHealthService(ctx, userRepo, userNatsPub)
	userService := initUserService(ctx, cfg, userRepo, userNatsPub)

	s, err := server.NewServer(
		server.SetConfig(cfg),
		server.SetLogger(logger),
		server.SetHandler(
			map[string]http.Handler{
				"":     health.MakeHTTPHandler(ctx, healthService),
				"user": user.MakeHTTPHandler(ctx, userService),
			}),
		server.SetGRPC(
			user.JoinGRPC(ctx, userService),
		),
	)
	if err != nil {
		level.Error(logger).Log("init", "server", "err", err)
		os.Exit(1)
	}
	defer s.Close()

	if err := s.AddHTTP(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	if err = s.AddGRPC(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	if err = s.AddMetrics(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	s.AddSignalHandler()
	s.Run()
}

func initHealthService(ctx context.Context, userRep userRepository.Repository, userNatsPub userQueue.Publisher) health.Service {
	var healthService health.Service
	healthService = health.NewHealthService(userRep, userNatsPub)
	healthService = health.NewLoggingService(ctx, healthService)
	return healthService
}

func initUserService(ctx context.Context, cfg *configs.Config, repo userRepository.Repository, ncPub userQueue.Publisher) user.Service {
	userService := user.NewUserService(repo, ncPub)
	if cfg.Metrics.Enabled {
		userService = user.NewMetricsService(ctx, userService)
	}
	userService = user.NewLoggingService(ctx, userService)
	if cfg.Tracer.Enabled {
		userService = user.NewTracingService(ctx, userService)
	}
	if cfg.Sentry.Enabled {
		userService = user.NewSentryService(userService)
	}
	return userService
}

func initUserRepository(ctx context.Context, db *database.Connection, cfg *configs.Config) userRepository.Repository {
	repo := userRepository.NewRepository(db)
	if cfg.Tracer.Enabled {
		repo = userRepository.NewTracingRepository(ctx, repo)
	}
	if cfg.Sentry.Enabled {
		repo = userRepository.NewSentryService(repo)
	}
	return repo
}
