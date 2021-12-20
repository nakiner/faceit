package user

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidArgument is returned when one or more arguments are invalid.
	ErrInvalidArgument = errors.New("invalid argument")
	ErrAlreadyExists   = errors.New("already exists")
	ErrBadRequest      = errors.New("bad request")
	ErrNotFound        = errors.New("not found")
	errBadRoute        = errors.New("bad route")
	ErrInvalidRequest  = errors.New("invalid params in request")
)

type ContextHTTPKey struct{}

type HTTPInfo struct {
	Method   string
	URL      string
	From     string
	Protocol string
}

type errorCode interface {
	Code() int
}

// getHTTPStatusCode returns http status code from error.
func getHTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if e, ok := err.(errorCode); ok && e.Code() != 0 {
		return e.Code()
	}

	switch errors.Cause(err) {
	case ErrInvalidArgument:
		return http.StatusBadRequest
	case ErrAlreadyExists:
		return http.StatusBadRequest
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
