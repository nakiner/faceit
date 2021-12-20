package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/nakiner/faceit/tools/logging"
	"github.com/nakiner/faceit/tools/tracing"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func MakeHTTPHandler(ctx context.Context, s Service) http.Handler {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "http handler", "user")
	tracer := tracing.FromContext(ctx)

	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		// httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(httpToContext()),
		httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "http server", logger)),
		httptransport.ServerFinalizer(closeHTTPTracer()),
	}

	r.Methods("POST").Path("/user").Handler(httptransport.NewServer(
		makeCreateUserEndpoint(s),
		decodePOSTCreateUserRequest,
		encodeCreateUserResponse,
		options...,
	))

	r.Methods("GET").Path("/user").Handler(httptransport.NewServer(
		makeGetUsersEndpoint(s),
		decodeGETGetUsersRequest,
		encodeGetUsersResponse,
		options...,
	))

	r.Methods("PUT").Path("/user/{id}").Handler(httptransport.NewServer(
		makeUpdateUserEndpoint(s),
		decodePUTUser,
		encodeStatus,
		options...,
	))

	r.Methods("DELETE").Path("/user/{id}").Handler(httptransport.NewServer(
		makeDeleteUserEndpoint(s),
		decodeDELETEDeleteUserRequest,
		encodeStatus,
		options...,
	))

	return accessControl(r)
}

func httpToContext() httptransport.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		return context.WithValue(ctx, ContextHTTPKey{}, HTTPInfo{
			Method:   req.Method,
			URL:      req.RequestURI,
			From:     req.RemoteAddr,
			Protocol: req.Proto,
		})
	}
}
func closeHTTPTracer() httptransport.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		span := stdopentracing.SpanFromContext(ctx)
		span.Finish()
	}
}

func decodePOSTCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}
	{
		if err := validate(request); err != nil {
			return nil, errors.Wrap(ErrInvalidRequest, err.Error())
		}
	}
	return request, nil
}

func decodeDELETEDeleteUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request DeleteUserRequest
	vars := mux.Vars(r)

	{

		id, ok := vars["id"]

		if !ok {
			return nil, errors.WithStack(errBadRoute)
		}

		request.Id = id
	}

	{
		if err := validate(request); err != nil {
			return nil, errors.Wrap(ErrInvalidRequest, err.Error())
		}
	}
	return request, nil
}

func decodeGETGetUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request GetUsersRequest

	{
		decoder := schema.NewDecoder()
		err := decoder.Decode(&request, r.URL.Query())
		if err != nil {
			return nil, errors.Wrap(ErrInvalidArgument, err.Error())
		}
	}
	{
		if err := validate(request); err != nil {
			return nil, errors.Wrap(ErrInvalidRequest, err.Error())
		}
	}
	return request, nil
}

func decodePUTUser(_ context.Context, r *http.Request) (interface{}, error) {
	var request User

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}

	vars := mux.Vars(r)

	{

		id, ok := vars["id"]

		if !ok {
			return nil, errors.WithStack(errBadRoute)
		}

		request.Id = id
	}

	{
		if err := validate(request); err != nil {
			return nil, errors.Wrap(ErrInvalidRequest, err.Error())
		}
	}
	return request, nil
}

func encodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeGetUsersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeStatus(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encodeError handles error from business-layer.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("X-Esp-Error", err.Error())
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")

	w.WriteHeader(getHTTPStatusCode(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// accessControl is CORS middleware.
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
