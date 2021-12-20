package health

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/nakiner/faceit/tools/logging"
	"github.com/pkg/errors"
)

func MakeHTTPHandler(ctx context.Context, s Service) http.Handler {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "http handler", "health")

	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		// httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(httpToContext()),
	}

	r.Methods("GET").Path("/liveness").Handler(httptransport.NewServer(
		makeLivenessEndpoint(s),
		decodeGETLivenessRequest,
		encodeLivenessResponse,
		options...,
	))

	r.Methods("GET").Path("/readiness").Handler(httptransport.NewServer(
		makeReadinessEndpoint(s),
		decodeGETReadinessRequest,
		encodeReadinessResponse,
		options...,
	))

	r.Methods("GET").Path("/version").Handler(httptransport.NewServer(
		makeVersionEndpoint(s),
		decodeGETVersionRequest,
		encodeVersionResponse,
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

func decodeGETLivenessRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request LivenessRequest

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

func decodeGETReadinessRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ReadinessRequest

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

func decodeGETVersionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request VersionRequest

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

func encodeLivenessResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeReadinessResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeVersionResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
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
