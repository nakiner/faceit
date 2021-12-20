package user

import "context"

type Repository interface {
	IsReady() bool
	Create(ctx context.Context, data *User) (string, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *User) error
	Get(ctx context.Context, conditions Conditions, limit uint32, offset uint32) ([]*User, error)
}
