package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport/grpc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	pb "github.com/nakiner/faceit/internal/faceitpb"
	"github.com/nakiner/faceit/tools/logging"
	"github.com/nakiner/faceit/tools/tracing"
	stdopentracing "github.com/opentracing/opentracing-go"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	createUser grpctransport.Handler
	getUsers   grpctransport.Handler
	updateUser grpctransport.Handler
	deleteUser grpctransport.Handler
}

type ContextGRPCKey struct{}

type GRPCInfo struct{}

// NewGRPCServer makes a set of endpoints available as a gRPC userServer.
func NewGRPCServer(ctx context.Context, s Service) pb.UserServiceServer {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "grpc handler", "user")
	tracer := tracing.FromContext(ctx)

	options := []grpctransport.ServerOption{
		// grpctransport.ServerErrorLogger(logger),
		grpctransport.ServerBefore(grpcToContext()),
		grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "grpc server", logger)),
		grpctransport.ServerFinalizer(closeGRPCTracer()),
	}

	return &grpcServer{
		createUser: grpctransport.NewServer(
			makeCreateUserEndpoint(s),
			decodeGRPCCreateUserRequest,
			encodeGRPCCreateUserResponse,
			options...,
		),
		getUsers: grpctransport.NewServer(
			makeGetUsersEndpoint(s),
			decodeGRPCGetUsersRequest,
			encodeGRPCGetUsersResponse,
			options...,
		),
		updateUser: grpctransport.NewServer(
			makeUpdateUserEndpoint(s),
			decodeGRPCUpdateUserRequest,
			encodeGRPCUpdateUserResponse,
			options...,
		),
		deleteUser: grpctransport.NewServer(
			makeDeleteUserEndpoint(s),
			decodeGRPCDeleteUserRequest,
			encodeGRPCDeleteUserResponse,
			options...,
		),
	}
}

func JoinGRPC(ctx context.Context, s Service) func(*googlegrpc.Server) {
	return func(g *googlegrpc.Server) {
		pb.RegisterUserServiceServer(g, NewGRPCServer(ctx, s))
	}
}

func grpcToContext() grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		return context.WithValue(ctx, ContextGRPCKey{}, GRPCInfo{})
	}
}
func closeGRPCTracer() grpc.ServerFinalizerFunc {
	return func(ctx context.Context, err error) {
		span := stdopentracing.SpanFromContext(ctx)
		span.Finish()
	}
}

func (s *grpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, rep, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateUserResponse), nil
}

func (s *grpcServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	_, rep, err := s.getUsers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUsersResponse), nil
}

func (s *grpcServer) UpdateUser(ctx context.Context, req *pb.User) (*pb.Status, error) {
	_, rep, err := s.updateUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.Status), nil
}

func (s *grpcServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Status, error) {
	_, rep, err := s.deleteUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.Status), nil
}

func decodeGRPCCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.CreateUserRequest)
	if !ok {
		return nil, errors.New("decodeGRPCCreateUserRequest wrong request")
	}

	req := PBToCreateUserRequest(inReq)
	if err := validate(req); err != nil {
		return nil, err
	}
	return *req, nil
}

func decodeGRPCGetUsersRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.GetUsersRequest)
	if !ok {
		return nil, errors.New("decodeGRPCGetUsersRequest wrong request")
	}

	req := PBToGetUsersRequest(inReq)
	if err := validate(req); err != nil {
		return nil, err
	}
	return *req, nil
}

func decodeGRPCUpdateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.User)
	if !ok {
		return nil, errors.New("decodeGRPCUpdateUserRequest wrong request")
	}

	req := PBToUser(inReq)
	if err := validate(req); err != nil {
		return nil, err
	}
	return *req, nil
}

func decodeGRPCDeleteUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.DeleteUserRequest)
	if !ok {
		return nil, errors.New("decodeGRPCDeleteUserRequest wrong request")
	}

	req := PBToDeleteUserRequest(inReq)
	if err := validate(req); err != nil {
		return nil, err
	}
	return *req, nil
}

func encodeGRPCCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*CreateUserResponse)
	if !ok {
		return nil, errors.New("encodeGRPCCreateUserResponse wrong response")
	}

	return CreateUserResponseToPB(inResp), nil
}

func encodeGRPCGetUsersResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*GetUsersResponse)
	if !ok {
		return nil, errors.New("encodeGRPCGetUsersResponse wrong response")
	}

	return GetUsersResponseToPB(inResp), nil
}

func encodeGRPCUpdateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*Status)
	if !ok {
		return nil, errors.New("encodeGRPCUpdateUserResponse wrong response")
	}

	return StatusToPB(inResp), nil
}

func encodeGRPCDeleteUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*Status)
	if !ok {
		return nil, errors.New("encodeGRPCDeleteUserResponse wrong response")
	}

	return StatusToPB(inResp), nil
}

func CreateUserRequestToPB(d *CreateUserRequest) *pb.CreateUserRequest {
	if d == nil {
		return nil
	}

	resp := pb.CreateUserRequest{
		FirstName:       d.FirstName,
		LastName:        d.LastName,
		Nickname:        d.Nickname,
		Password:        d.Password,
		PasswordConfirm: d.PasswordConfirm,
		Email:           d.Email,
		Country:         d.Country,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}

	return &resp
}

func PBToCreateUserRequest(d *pb.CreateUserRequest) *CreateUserRequest {
	if d == nil {
		return nil
	}

	resp := CreateUserRequest{
		FirstName:       d.FirstName,
		LastName:        d.LastName,
		Nickname:        d.Nickname,
		Password:        d.Password,
		PasswordConfirm: d.PasswordConfirm,
		Email:           d.Email,
		Country:         d.Country,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}

	return &resp
}

func CreateUserResponseToPB(d *CreateUserResponse) *pb.CreateUserResponse {
	if d == nil {
		return nil
	}

	resp := pb.CreateUserResponse{
		Id: d.Id,
	}

	return &resp
}

func PBToCreateUserResponse(d *pb.CreateUserResponse) *CreateUserResponse {
	if d == nil {
		return nil
	}

	resp := CreateUserResponse{
		Id: d.Id,
	}

	return &resp
}

func DeleteUserRequestToPB(d *DeleteUserRequest) *pb.DeleteUserRequest {
	if d == nil {
		return nil
	}

	resp := pb.DeleteUserRequest{
		Id: d.Id,
	}

	return &resp
}

func PBToDeleteUserRequest(d *pb.DeleteUserRequest) *DeleteUserRequest {
	if d == nil {
		return nil
	}

	resp := DeleteUserRequest{
		Id: d.Id,
	}

	return &resp
}

func GetUsersRequestToPB(d *GetUsersRequest) *pb.GetUsersRequest {
	if d == nil {
		return nil
	}

	resp := pb.GetUsersRequest{
		Limit:     d.Limit,
		Offset:    d.Offset,
		Id:        d.Id,
		Country:   d.Country,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Nickname:  d.Nickname,
	}

	return &resp
}

func PBToGetUsersRequest(d *pb.GetUsersRequest) *GetUsersRequest {
	if d == nil {
		return nil
	}

	resp := GetUsersRequest{
		Limit:     d.Limit,
		Offset:    d.Offset,
		Id:        d.Id,
		Country:   d.Country,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Nickname:  d.Nickname,
	}

	return &resp
}

func GetUsersResponseToPB(d *GetUsersResponse) *pb.GetUsersResponse {
	if d == nil {
		return nil
	}

	resp := pb.GetUsersResponse{}

	for _, v := range *d {
		resp.Data = append(resp.Data, UserToPB(&v))
	}

	return &resp
}

func PBToGetUsersResponse(d *pb.GetUsersResponse) *GetUsersResponse {
	if d == nil {
		return nil
	}

	resp := GetUsersResponse{}

	for _, v := range d.Data {
		resp = append(resp, *PBToUser(v))
	}

	return &resp
}

func StatusToPB(d *Status) *pb.Status {
	if d == nil {
		return nil
	}

	resp := pb.Status{
		Status:  d.Status,
		Message: d.Message,
	}

	return &resp
}

func PBToStatus(d *pb.Status) *Status {
	if d == nil {
		return nil
	}

	resp := Status{
		Status:  d.Status,
		Message: d.Message,
	}

	return &resp
}

func UserToPB(d *User) *pb.User {
	if d == nil {
		return nil
	}

	resp := pb.User{
		Id:        d.Id,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Nickname:  d.Nickname,
		Password:  d.Password,
		Email:     d.Email,
		Country:   d.Country,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}

	return &resp
}

func PBToUser(d *pb.User) *User {
	if d == nil {
		return nil
	}

	resp := User{
		Id:        d.Id,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Nickname:  d.Nickname,
		Password:  d.Password,
		Email:     d.Email,
		Country:   d.Country,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}

	return &resp
}
