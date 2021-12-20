package configs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	assert.Equal(t, &Config{}, NewConfig())
}

func TestConfig_Read(t *testing.T) {
	port := 9154
	exp := &Config{}
	exp.Metrics.Port = port
	os.Setenv("FACEIT_METRICS_PORT", fmt.Sprintf("%d", port))
	c := NewConfig()
	c.Read()
	assert.Equal(t, exp.Metrics.Port, c.Metrics.Port)
}
