//go:generate easyjson -all endpoint.go
package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	_ "github.com/mailru/easyjson/gen"
)

//easyjson:json
type CreateUserRequest struct {
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	Nickname        string `json:"nickname,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
	Email           string `json:"email,omitempty"`
	Country         string `json:"country,omitempty"`
	CreatedAt       string `json:"createdAt,omitempty"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
}

//easyjson:json
type CreateUserResponse struct {
	Id string `json:"id,omitempty"`
}

//easyjson:json
type DeleteUserRequest struct {
	Id string `json:"id,omitempty"`
}

type GetUsersRequest struct {
	Limit     uint32 `json:"limit,omitempty"`
	Offset    uint32 `json:"offset,omitempty"`
	Id        string `schema:"id"`
	Country   string `schema:"country"`
	FirstName string `schema:"firstName"`
	LastName  string `schema:"lastName"`
	Nickname  string `schema:"nickname"`
}

//easyjson:json
type GetUsersResponse []User

//easyjson:json
type Status struct {
	Status  bool   `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

//easyjson:json
type User struct {
	Id        string `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email,omitempty"`
	Country   string `json:"country,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

//easyjson:skip
type endpoints struct {
	CreateUserEndpoint endpoint.Endpoint
	GetUsersEndpoint   endpoint.Endpoint
	UpdateUserEndpoint endpoint.Endpoint
	DeleteUserEndpoint endpoint.Endpoint
}

func (e endpoints) CreateUser(ctx context.Context, req *CreateUserRequest) (resp *CreateUserResponse, err error) {
	response, err := e.CreateUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(CreateUserResponse)
	return &r, err
}

func (e endpoints) GetUsers(ctx context.Context, req *GetUsersRequest) (resp *GetUsersResponse, err error) {
	response, err := e.GetUsersEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(GetUsersResponse)
	return &r, err
}

func (e endpoints) UpdateUser(ctx context.Context, req *User) (resp *Status, err error) {
	response, err := e.UpdateUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(Status)
	return &r, err
}

func (e endpoints) DeleteUser(ctx context.Context, req *DeleteUserRequest) (resp *Status, err error) {
	response, err := e.DeleteUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(Status)
	return &r, err
}

func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		return s.CreateUser(ctx, &req)
	}
}

func makeGetUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUsersRequest)
		return s.GetUsers(ctx, &req)
	}
}

func makeUpdateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(User)
		return s.UpdateUser(ctx, &req)
	}
}

func makeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteUserRequest)
		return s.DeleteUser(ctx, &req)
	}
}
