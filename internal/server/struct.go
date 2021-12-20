package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nakiner/faceit/configs"
	"github.com/nakiner/faceit/tools/limiting"
	"github.com/nakiner/faceit/tools/sentry"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// Server main struct for prm-export service
type Server struct {
	cfg     *configs.Config
	logger  log.Logger
	handler http.Handler
	grpc    *grpc.Server
	group   run.Group
}

type Option func(*Server)

func (s *Server) setGroup(group run.Group) {
	s.group = group
}

// Run запускает сервер
func (s *Server) Run() error {
	return s.logger.Log("exit", s.group.Run())
}

// Close closes everything that is open and should be closed after server shutdown
func (s *Server) Close() {

}

// AddHTTP  http server start when Server.Run()
func (s *Server) AddHTTP() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.HTTP.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "cann't add HTTP transport")
	}

	s.group.Add(func() error {
		level.Info(s.logger).Log("component", "HTTP server", "addr", addr, "msg", "listening...")
		if s.cfg.Limiter.Enabled {
			l := limiting.NewLimiter(context.Background(), s.cfg.Limiter.Limit)
			s.handler = l.Middleware(s.handler)
		}
		if s.cfg.Sentry.Enabled {
			s.handler = sentry.Middleware(s.handler)
		}

		httpServer := &http.Server{
			Handler:      accessControl(s.handler),
			WriteTimeout: time.Second * time.Duration(s.cfg.Server.HTTP.TimeoutSec),
		}

		return httpServer.Serve(listener)
	}, func(error) {
		listener.Close()
	})
	return nil
}

// AddGRPC  grpc server start when Server.Run()
func (s *Server) AddGRPC() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "cann't add GRPC transport")
	}

	s.group.Add(func() error {
		level.Info(s.logger).Log("component", "GRPC server", "addr", addr, "msg", "listening...")
		return s.grpc.Serve(listener)
	}, func(error) {
		listener.Close()
	})
	return nil
}

// AddMetrics  metrics server start when Server.Run()
func (s *Server) AddMetrics() error {
	if !s.cfg.Metrics.Enabled {
		return nil
	}
	addr := fmt.Sprintf(":%d", s.cfg.Metrics.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "cann't add Metrics")
	}
	s.group.Add(func() error {
		level.Info(s.logger).Log("component", "metrics server", "addr", addr, "msg", "listening...")
		http.Handle("/metrics", promhttp.Handler())
		return http.Serve(listener, http.DefaultServeMux)
	}, func(error) {
		listener.Close()
	})

	return nil
}

// AddSignalHandler add listener os signal when Server.Run()
func (s *Server) AddSignalHandler() {
	ch := make(chan struct{})
	s.group.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return errors.Errorf("received signal %s", sig)
		case <-ch:
			return nil
		}
	}, func(error) {
		close(ch)
	})
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
