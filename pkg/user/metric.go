package user

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	tool "github.com/nakiner/faceit/tools/metrics"
)

// NewMetricService returns an instance of an instrumenting Service.
func NewMetricsService(ctx context.Context, s Service) Service {
	counter, latency := tool.FromContext(ctx)
	return &metricService{counter, latency, s}
}

type metricService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func (s *metricService) CreateUser(ctx context.Context, req *CreateUserRequest) (resp *CreateUserResponse, err error) {
	defer func(begin time.Time) {
		go func() {
			s.requestCount.With("service", "user", "handler", "CreateUser", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
			s.requestLatency.With("service", "user", "handler", "CreateUser", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
		}()
	}(time.Now())
	return s.Service.CreateUser(ctx, req)
}

func (s *metricService) GetUsers(ctx context.Context, req *GetUsersRequest) (resp *GetUsersResponse, err error) {
	defer func(begin time.Time) {
		go func() {
			s.requestCount.With("service", "user", "handler", "GetUsers", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
			s.requestLatency.With("service", "user", "handler", "GetUsers", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
		}()
	}(time.Now())
	return s.Service.GetUsers(ctx, req)
}

func (s *metricService) UpdateUser(ctx context.Context, req *User) (resp *Status, err error) {
	defer func(begin time.Time) {
		go func() {
			s.requestCount.With("service", "user", "handler", "UpdateUser", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
			s.requestLatency.With("service", "user", "handler", "UpdateUser", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
		}()
	}(time.Now())
	return s.Service.UpdateUser(ctx, req)
}

func (s *metricService) DeleteUser(ctx context.Context, req *DeleteUserRequest) (resp *Status, err error) {
	defer func(begin time.Time) {
		go func() {
			s.requestCount.With("service", "user", "handler", "DeleteUser", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
			s.requestLatency.With("service", "user", "handler", "DeleteUser", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
		}()
	}(time.Now())
	return s.Service.DeleteUser(ctx, req)
}
