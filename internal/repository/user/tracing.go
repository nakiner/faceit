package user

import (
	"context"
	"github.com/nakiner/faceit/tools/tracing"
	"github.com/opentracing/opentracing-go"
)

// NewTracingRepository allows overriding behavior on given repository and traces into opentracing handler
// allows sending detailed traces with preset params and debug data
func NewTracingRepository(ctx context.Context, r Repository) Repository {
	tracer := tracing.FromContext(ctx)
	return &tracingRepository{tracer, r}
}

type tracingRepository struct {
	tracer opentracing.Tracer
	Repository
}

func (r *tracingRepository) IsReady() bool {
	return r.Repository.IsReady()
}

func (r *tracingRepository) Create(ctx context.Context, data *User) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "Create")
	defer span.Finish()
	return r.Repository.Create(ctx, data)
}

func (r *tracingRepository) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "Delete")
	defer span.Finish()
	return r.Repository.Delete(ctx, id)
}

func (r *tracingRepository) Update(ctx context.Context, data *User) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "Update")
	defer span.Finish()
	return r.Repository.Update(ctx, data)
}

func (r *tracingRepository) Get(ctx context.Context, conditions Conditions, limit uint32, offset uint32) ([]*User, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "Delete")
	defer span.Finish()
	return r.Repository.Get(ctx, conditions, limit, offset)
}
