//go:generate easyjson -all endpoint.go
package health

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	_ "github.com/mailru/easyjson/gen"
)

//easyjson:json
type LivenessRequest struct {
}

//easyjson:json
type LivenessResponse struct {
	Status string `json:"status,omitempty"`
}

//easyjson:json
type ReadinessRequest struct {
}

//easyjson:json
type ReadinessResponse struct {
	Status string `json:"status,omitempty"`
}

//easyjson:json
type VersionRequest struct {
}

//easyjson:json
type VersionResponse struct {
	BuildTime string `json:"BuildTime,omitempty"`
	Version   string `json:"Version,omitempty"`
	Commit    string `json:"Commit,omitempty"`
}

//easyjson:skip
type endpoints struct {
	LivenessEndpoint  endpoint.Endpoint
	ReadinessEndpoint endpoint.Endpoint
	VersionEndpoint   endpoint.Endpoint
}

func (e endpoints) Liveness(ctx context.Context, req *LivenessRequest) (resp *LivenessResponse, err error) {
	response, err := e.LivenessEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(LivenessResponse)
	return &r, err
}

func (e endpoints) Readiness(ctx context.Context, req *ReadinessRequest) (resp *ReadinessResponse, err error) {
	response, err := e.ReadinessEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(ReadinessResponse)
	return &r, err
}

func (e endpoints) Version(ctx context.Context, req *VersionRequest) (resp *VersionResponse, err error) {
	response, err := e.VersionEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(VersionResponse)
	return &r, err
}

func makeLivenessEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LivenessRequest)
		return s.Liveness(ctx, &req)
	}
}

func makeReadinessEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ReadinessRequest)
		return s.Readiness(ctx, &req)
	}
}

func makeVersionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VersionRequest)
		return s.Version(ctx, &req)
	}
}
