package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// NewHTTPClient returns an Service backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string, tracer stdopentracing.Tracer, logger log.Logger) (Service, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// global client middlewares
	var options []httptransport.ClientOption
	if tracer != nil {
		options = append(
			options,
			httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)),
		)
	}

	return endpoints{
		CreateUserEndpoint: httptransport.NewClient(
			"POST",
			copyURL(u, "/user"),
			encodeHTTPCreateUserCreateUserRequest,
			decodeHTTPCreateUserCreateUserResponse,
			options...,
		).Endpoint(),
		GetUsersEndpoint: httptransport.NewClient(
			"GET",
			copyURL(u, "/user"),
			encodeHTTPGetUsersGetUsersRequest,
			decodeHTTPGetUsersGetUsersResponse,
			options...,
		).Endpoint(),
		UpdateUserEndpoint: httptransport.NewClient(
			"PUT",
			copyURL(u, "/user/{id}"),
			encodeHTTPUpdateUserUser,
			decodeHTTPUpdateUserStatus,
			options...,
		).Endpoint(),
		DeleteUserEndpoint: httptransport.NewClient(
			"DELETE",
			copyURL(u, "/user/{id}"),
			encodeHTTPDeleteUserDeleteUserRequest,
			decodeHTTPDeleteUserStatus,
			options...,
		).Endpoint(),
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func encodeHTTPCreateUserCreateUserRequest(_ context.Context, r *http.Request, request interface{}) error {

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return errors.Wrap(err, "encode request body")
	}
	r.Body = ioutil.NopCloser(&buf)

	return nil
}
func encodeHTTPGetUsersGetUsersRequest(_ context.Context, r *http.Request, request interface{}) error {
	{
		queryMap := make(map[string][]string)
		if err := schema.NewEncoder().Encode(request, queryMap); err == nil {
			query := url.Values(queryMap)
			r.URL.RawQuery = query.Encode()
		}
	}

	return nil
}
func encodeHTTPUpdateUserUser(_ context.Context, r *http.Request, request interface{}) error {

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return errors.Wrap(err, "encode request body")
	}
	r.Body = ioutil.NopCloser(&buf)
	req := request.(*User)
	rout := mux.NewRouter()
	rout.Path(r.URL.Path).Name("UpdateUser")

	url, err := rout.Get("UpdateUser").URL(
		"id", fmt.Sprint(req.Id),
	)
	if err != nil {
		return err
	}

	r.URL.Path = url.String()

	return nil
}
func encodeHTTPDeleteUserDeleteUserRequest(_ context.Context, r *http.Request, request interface{}) error {
	{
		queryMap := make(map[string][]string)
		if err := schema.NewEncoder().Encode(request, queryMap); err == nil {
			query := url.Values(queryMap)
			r.URL.RawQuery = query.Encode()
		}
	}
	req := request.(*DeleteUserRequest)
	rout := mux.NewRouter()
	rout.Path(r.URL.Path).Name("DeleteUser")

	url, err := rout.Get("DeleteUser").URL(
		"id", fmt.Sprint(req.Id),
	)
	if err != nil {
		return err
	}

	r.URL.Path = url.String()

	return nil
}

func decodeHTTPCreateUserCreateUserResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var request CreateUserResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}
	return request, nil
}

func decodeHTTPGetUsersGetUsersResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var request GetUsersResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}
	return request, nil
}

func decodeHTTPUpdateUserStatus(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var request Status
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}
	return request, nil
}

func decodeHTTPDeleteUserStatus(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var request Status
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}
	return request, nil
}
