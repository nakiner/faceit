//go:build integration && !unit
// +build integration,!unit

package integration

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/nakiner/faceit/pkg/health"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

const htttAddrhealth = "localhost:8080"

func TestHTTPHealthServiceLiveness(t *testing.T) {
	client, err := health.NewHTTPClient(htttAddrhealth, opentracing.GlobalTracer(), log.NewNopLogger())
	assert.NoError(t, err)
	_, err = client.Liveness(context.Background(), &health.LivenessRequest{})
	assert.NoError(t, err)
}

func TestHTTPHealthServiceReadiness(t *testing.T) {
	client, err := health.NewHTTPClient(htttAddrhealth, opentracing.GlobalTracer(), log.NewNopLogger())
	assert.NoError(t, err)
	_, err = client.Readiness(context.Background(), &health.ReadinessRequest{})
	assert.NoError(t, err)
}

func TestHTTPHealthServiceVersion(t *testing.T) {
	client, err := health.NewHTTPClient(htttAddrhealth, opentracing.GlobalTracer(), log.NewNopLogger())
	assert.NoError(t, err)
	_, err = client.Version(context.Background(), &health.VersionRequest{})
	assert.NoError(t, err)
}
