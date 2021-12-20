package user

import (
	"context"
	"strconv"

	"github.com/getsentry/sentry-go"
)

func NewSentryService(s Service) Service {
	return &sentryService{s}
}

type sentryService struct {
	Service
}

type sentryLog interface {
	SentryLog() []interface{}
}

func (s *sentryService) getSentryLog(req interface{}, resp interface{}) (out map[string][]interface{}) {
	out = make(map[string][]interface{})
	if sentry, ok := interface{}(req).(sentryLog); ok {
		out["request"] = append(out["request"], sentry.SentryLog()...)
	}

	if sentry, ok := interface{}(resp).(sentryLog); ok {
		out["response"] = append(out["response"], sentry.SentryLog()...)
	}
	return
}

func (s *sentryService) CreateUser(ctx context.Context, req *CreateUserRequest) (resp *CreateUserResponse, err error) {
	defer func() {
		if err != nil {
			log := s.getSentryLog(req, resp)
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", strconv.Itoa(getHTTPStatusCode(err)))
				scope.SetTag("method", "CreateUser")
				scope.SetExtra("request", log["request"])
				scope.SetExtra("response", log["response"])
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Service.CreateUser(ctx, req)
}

func (s *sentryService) GetUsers(ctx context.Context, req *GetUsersRequest) (resp *GetUsersResponse, err error) {
	defer func() {
		if err != nil {
			log := s.getSentryLog(req, resp)
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", strconv.Itoa(getHTTPStatusCode(err)))
				scope.SetTag("method", "GetUsers")
				scope.SetExtra("request", log["request"])
				scope.SetExtra("response", log["response"])
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Service.GetUsers(ctx, req)
}

func (s *sentryService) UpdateUser(ctx context.Context, req *User) (resp *Status, err error) {
	defer func() {
		if err != nil {
			log := s.getSentryLog(req, resp)
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", strconv.Itoa(getHTTPStatusCode(err)))
				scope.SetTag("method", "UpdateUser")
				scope.SetExtra("request", log["request"])
				scope.SetExtra("response", log["response"])
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Service.UpdateUser(ctx, req)
}

func (s *sentryService) DeleteUser(ctx context.Context, req *DeleteUserRequest) (resp *Status, err error) {
	defer func() {
		if err != nil {
			log := s.getSentryLog(req, resp)
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", strconv.Itoa(getHTTPStatusCode(err)))
				scope.SetTag("method", "DeleteUser")
				scope.SetExtra("request", log["request"])
				scope.SetExtra("response", log["response"])
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Service.DeleteUser(ctx, req)
}
