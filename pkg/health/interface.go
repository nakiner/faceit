//go:generate mockgen -destination service_mock.go -package health  github.com/nakiner/faceit/pkg/health Service
package health

import (
	"context"

	_ "github.com/golang/mock/mockgen/model"
)

type Service interface {

	// Liveness returns a error if service doesn`t live.
	Liveness(context.Context, *LivenessRequest) (*LivenessResponse, error)

	// Readiness returns a error if service doesn`t ready.
	Readiness(context.Context, *ReadinessRequest) (*ReadinessResponse, error)

	// Version returns build time, last commit and version app
	Version(context.Context, *VersionRequest) (*VersionResponse, error)
}
