package user

import (
	"context"

	"github.com/nakiner/faceit/tools/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// NewTracingService returns an instance of an instrumenting Service.
func NewTracingService(ctx context.Context, s Service) Service {
	tracer := tracing.FromContext(ctx)
	return &tracingService{tracer, s}
}

type trace interface {
	Span() []interface{}
}

func (s *tracingService) getTrace(req interface{}, resp interface{}) (out []interface{}) {
	if val, ok := interface{}(req).(trace); ok {
		out = append(out, val.Span()...)
	}

	if val, ok := interface{}(resp).(trace); ok {
		out = append(out, val.Span()...)
	}

	return
}

type tracingService struct {
	tracer opentracing.Tracer
	Service
}

func (s *tracingService) CreateUser(ctx context.Context, req *CreateUserRequest) (resp *CreateUserResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateUser")
	span.LogFields(log.Object("tracingService", s.getTrace(req, resp)))
	defer span.Finish()
	return s.Service.CreateUser(ctx, req)
}

func (s *tracingService) GetUsers(ctx context.Context, req *GetUsersRequest) (resp *GetUsersResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetUsers")
	span.LogFields(log.Object("tracingService", s.getTrace(req, resp)))
	defer span.Finish()
	return s.Service.GetUsers(ctx, req)
}

func (s *tracingService) UpdateUser(ctx context.Context, req *User) (resp *Status, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateUser")
	span.LogFields(log.Object("tracingService", s.getTrace(req, resp)))
	defer span.Finish()
	return s.Service.UpdateUser(ctx, req)
}

func (s *tracingService) DeleteUser(ctx context.Context, req *DeleteUserRequest) (resp *Status, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteUser")
	span.LogFields(log.Object("tracingService", s.getTrace(req, resp)))
	defer span.Finish()
	return s.Service.DeleteUser(ctx, req)
}
