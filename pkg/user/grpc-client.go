package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	pb "github.com/nakiner/faceit/internal/faceitpb"
	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// NewGRPCClient returns an Service backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) Service {
	// global client middlewares
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
	}

	return endpoints{
		// Each individual endpoint is an grpc/transport.Client (which implements
		// endpoint.Endpoint) that gets wrapped with various middlewares. If you
		// made your own client library, you'd do this work there, so your server
		// could rely on a consistent set of client behavior.
		CreateUserEndpoint: grpctransport.NewClient(
			conn,
			"faceitpb.UserService",
			"CreateUser",
			encodeGRPCCreateUserRequest,
			decodeGRPCCreateUserResponse,
			pb.CreateUserResponse{},
			options...,
		).Endpoint(),
		GetUsersEndpoint: grpctransport.NewClient(
			conn,
			"faceitpb.UserService",
			"GetUsers",
			encodeGRPCGetUsersRequest,
			decodeGRPCGetUsersResponse,
			pb.GetUsersResponse{},
			options...,
		).Endpoint(),
		UpdateUserEndpoint: grpctransport.NewClient(
			conn,
			"faceitpb.UserService",
			"UpdateUser",
			encodeGRPCUser,
			decodeGRPCStatus,
			pb.Status{},
			options...,
		).Endpoint(),
		DeleteUserEndpoint: grpctransport.NewClient(
			conn,
			"faceitpb.UserService",
			"DeleteUser",
			encodeGRPCDeleteUserRequest,
			decodeGRPCStatus,
			pb.Status{},
			options...,
		).Endpoint(),
	}
}

func encodeGRPCCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*CreateUserRequest)
	if !ok {
		return nil, errors.New("encodeGRPCCreateUserRequest wrong request")
	}

	return CreateUserRequestToPB(inReq), nil
}

func encodeGRPCDeleteUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*DeleteUserRequest)
	if !ok {
		return nil, errors.New("encodeGRPCDeleteUserRequest wrong request")
	}

	return DeleteUserRequestToPB(inReq), nil
}

func encodeGRPCGetUsersRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*GetUsersRequest)
	if !ok {
		return nil, errors.New("encodeGRPCGetUsersRequest wrong request")
	}

	return GetUsersRequestToPB(inReq), nil
}

func encodeGRPCUser(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*User)
	if !ok {
		return nil, errors.New("encodeGRPCUser wrong request")
	}

	return UserToPB(inReq), nil
}

func decodeGRPCCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*pb.CreateUserResponse)
	if !ok {
		return nil, errors.New("decodeGRPCCreateUserResponse wrong response")
	}

	resp := PBToCreateUserResponse(inResp)

	return *resp, nil
}

func decodeGRPCGetUsersResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*pb.GetUsersResponse)
	if !ok {
		return nil, errors.New("decodeGRPCGetUsersResponse wrong response")
	}

	resp := PBToGetUsersResponse(inResp)

	return *resp, nil
}

func decodeGRPCStatus(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*pb.Status)
	if !ok {
		return nil, errors.New("decodeGRPCStatus wrong response")
	}

	resp := PBToStatus(inResp)

	return *resp, nil
}
