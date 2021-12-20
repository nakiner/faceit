//go:build integration && !unit
// +build integration,!unit

package integration

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/nakiner/faceit/pkg/user"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

const htttAddruser = "localhost:8080"

func TestHTTPUserServiceCreateUser(t *testing.T) {
	client, err := user.NewHTTPClient(htttAddruser, opentracing.GlobalTracer(), log.NewNopLogger())
	assert.NoError(t, err)
	_, err = client.CreateUser(context.Background(), &user.CreateUserRequest{Nickname: "sample"})
	assert.NoError(t, err)
}

func TestHTTPUserServiceGetUsers(t *testing.T) {
	client, err := user.NewHTTPClient(htttAddruser, opentracing.GlobalTracer(), log.NewNopLogger())
	assert.NoError(t, err)
	_, err = client.GetUsers(context.Background(), &user.GetUsersRequest{})
	assert.NoError(t, err)
}

func TestHTTPUserServiceUpdateUser(t *testing.T) {
	client, err := user.NewHTTPClient(htttAddruser, opentracing.GlobalTracer(), log.NewNopLogger())
	assert.NoError(t, err)
	resp, err := client.CreateUser(context.Background(), &user.CreateUserRequest{Nickname: "sample"})
	assert.NoError(t, err)
	_, err = client.UpdateUser(context.Background(), &user.User{Id: resp.Id, Nickname: "sample"})
	assert.NoError(t, err)
}

func TestHTTPUserServiceDeleteUser(t *testing.T) {
	client, err := user.NewHTTPClient(htttAddruser, opentracing.GlobalTracer(), log.NewNopLogger())
	assert.NoError(t, err)
	resp, err := client.CreateUser(context.Background(), &user.CreateUserRequest{Nickname: "sample"})
	assert.NoError(t, err)
	_, err = client.DeleteUser(context.Background(), &user.DeleteUserRequest{Id: resp.Id})
	assert.NoError(t, err)
}
