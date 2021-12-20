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
	"google.golang.org/grpc"
)

const grpcAddruser = "localhost:9194"

func TestGRPCUserServiceCreateUser(t *testing.T) {

	conn, err := grpc.Dial(grpcAddruser, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connection to grpc server: %s", err)
	}
	defer conn.Close()

	client := user.NewGRPCClient(conn, opentracing.GlobalTracer(), log.NewNopLogger())
	_, err = client.CreateUser(context.Background(), &user.CreateUserRequest{})

	assert.NoError(t, err)
}

func TestGRPCUserServiceGetUsers(t *testing.T) {

	conn, err := grpc.Dial(grpcAddruser, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connection to grpc server: %s", err)
	}
	defer conn.Close()

	client := user.NewGRPCClient(conn, opentracing.GlobalTracer(), log.NewNopLogger())
	_, err = client.GetUsers(context.Background(), &user.GetUsersRequest{})

	assert.NoError(t, err)
}

func TestGRPCUserServiceUpdateUser(t *testing.T) {

	conn, err := grpc.Dial(grpcAddruser, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connection to grpc server: %s", err)
	}
	defer conn.Close()

	client := user.NewGRPCClient(conn, opentracing.GlobalTracer(), log.NewNopLogger())
	_, err = client.UpdateUser(context.Background(), &user.User{})

	assert.NoError(t, err)
}

func TestGRPCUserServiceDeleteUser(t *testing.T) {

	conn, err := grpc.Dial(grpcAddruser, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connection to grpc server: %s", err)
	}
	defer conn.Close()

	client := user.NewGRPCClient(conn, opentracing.GlobalTracer(), log.NewNopLogger())
	_, err = client.DeleteUser(context.Background(), &user.DeleteUserRequest{})

	assert.NoError(t, err)
}
