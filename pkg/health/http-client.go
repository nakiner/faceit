package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
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
		LivenessEndpoint: httptransport.NewClient(
			"GET",
			copyURL(u, "/liveness"),
			encodeHTTPLivenessLivenessRequest,
			decodeHTTPLivenessLivenessResponse,
			options...,
		).Endpoint(),
		ReadinessEndpoint: httptransport.NewClient(
			"GET",
			copyURL(u, "/readiness"),
			encodeHTTPReadinessReadinessRequest,
			decodeHTTPReadinessReadinessResponse,
			options...,
		).Endpoint(),
		VersionEndpoint: httptransport.NewClient(
			"GET",
			copyURL(u, "/version"),
			encodeHTTPVersionVersionRequest,
			decodeHTTPVersionVersionResponse,
			options...,
		).Endpoint(),
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func encodeHTTPLivenessLivenessRequest(_ context.Context, r *http.Request, request interface{}) error {
	{
		queryMap := make(map[string][]string)
		if err := schema.NewEncoder().Encode(request, queryMap); err == nil {
			query := url.Values(queryMap)
			r.URL.RawQuery = query.Encode()
		}
	}

	return nil
}
func encodeHTTPReadinessReadinessRequest(_ context.Context, r *http.Request, request interface{}) error {
	{
		queryMap := make(map[string][]string)
		if err := schema.NewEncoder().Encode(request, queryMap); err == nil {
			query := url.Values(queryMap)
			r.URL.RawQuery = query.Encode()
		}
	}

	return nil
}
func encodeHTTPVersionVersionRequest(_ context.Context, r *http.Request, request interface{}) error {
	{
		queryMap := make(map[string][]string)
		if err := schema.NewEncoder().Encode(request, queryMap); err == nil {
			query := url.Values(queryMap)
			r.URL.RawQuery = query.Encode()
		}
	}

	return nil
}

func decodeHTTPLivenessLivenessResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var request LivenessResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}
	return request, nil
}

func decodeHTTPReadinessReadinessResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var request ReadinessResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}
	return request, nil
}

func decodeHTTPVersionVersionResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var request VersionResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}
	return request, nil
}
