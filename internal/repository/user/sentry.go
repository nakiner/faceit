package user

import (
	"context"
	"github.com/getsentry/sentry-go"
)

// NewSentryService allows overriding behavior on given repository and send events to sentry
func NewSentryService(r Repository) Repository {
	return &sentryRepository{r}
}

type sentryRepository struct {
	Repository
}

func (s *sentryRepository) IsReady() bool {
	return s.Repository.IsReady()
}

func (s *sentryRepository) Create(ctx context.Context, data *User) (id string, err error) {
	defer func() {
		if err != nil {
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("repository", "userDBRepository")
				scope.SetTag("method", "Create")
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Repository.Create(ctx, data)
}

func (s *sentryRepository) Delete(ctx context.Context, id string) (err error) {
	defer func() {
		if err != nil {
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("repository", "userDBRepository")
				scope.SetTag("method", "Delete")
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Repository.Delete(ctx, id)
}

func (s *sentryRepository) Update(ctx context.Context, data *User) (err error) {
	defer func() {
		if err != nil {
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("repository", "userDBRepository")
				scope.SetTag("method", "Update")
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Repository.Update(ctx, data)
}

func (s *sentryRepository) Get(ctx context.Context, conditions Conditions, limit uint32, offset uint32) (u []*User, err error) {
	defer func() {
		if err != nil {
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("repository", "userDBRepository")
				scope.SetTag("method", "Get")
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Repository.Get(ctx, conditions, limit, offset)
}
