package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type Config struct {
	Host           string
	Port           int
	UserName       string
	Password       string
	RequestTimeOut int `mapstructure:"request_timeout_msec"`
	RetryLimit     int `mapstructure:"retry_limit"`
	WaitLimit      int `mapstructure:"reconnect_time_wait_msec"`
}

func NewClient(cfg *Config) (*nats.Conn, error) {
	return nats.Connect(
		fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(cfg.RetryLimit),
		nats.ReconnectWait(time.Millisecond*time.Duration(cfg.WaitLimit)),
		nats.UserInfo(cfg.UserName, cfg.Password),
	)
}

func NewEncodedClient(nc *nats.Conn) (*nats.EncodedConn, error) {
	return nats.NewEncodedConn(nc, nats.JSON_ENCODER)
}
