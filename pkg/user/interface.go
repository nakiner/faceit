//go:generate mockgen -destination service_mock.go -package user  github.com/nakiner/faceit/pkg/user Service
package user

import (
	"context"

	_ "github.com/golang/mock/mockgen/model"
)

type Service interface {

	// CreateUser Create a new user
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)

	// GetUsers Get existing users, possibly allowing filter by arguments
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)

	// UpdateUser Update existing user
	UpdateUser(context.Context, *User) (*Status, error)

	// DeleteUser Delete existing user
	DeleteUser(context.Context, *DeleteUserRequest) (*Status, error)
}
