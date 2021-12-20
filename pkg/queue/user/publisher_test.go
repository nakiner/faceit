package user

import (
	"fmt"
	natsserver "github.com/nats-io/nats-server/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
)

var port = 4223

func TestNewPublisher(t *testing.T) {
	opt := natsserver.DefaultTestOptions
	opt.Port = port
	natsSvr := natsserver.RunServer(&opt)
	natsSvr.Start()
	defer natsSvr.Shutdown()

	nc, err := nats.Connect(
		fmt.Sprintf("nats://%s:%d", opt.Host, opt.Port),
	)
	assert.NoError(t, err)
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	assert.NoError(t, err)
	defer ec.Close()

	pub, err := NewPublisher(ec)
	assert.NoError(t, err)
	assert.NotNil(t, pub)
}

func TestPublisher_IsReady(t *testing.T) {
	opt := natsserver.DefaultTestOptions
	opt.Port = port
	natsSvr := natsserver.RunServer(&opt)
	natsSvr.Start()
	defer natsSvr.Shutdown()

	nc, err := nats.Connect(
		fmt.Sprintf("nats://%s:%d", opt.Host, opt.Port),
	)
	assert.NoError(t, err)
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	assert.NoError(t, err)
	defer ec.Close()

	pub, err := NewPublisher(ec)
	assert.NoError(t, err)
	assert.True(t, pub.IsReady())
}

func TestPublisher_UpdateUser(t *testing.T) {
	opt := natsserver.DefaultTestOptions
	opt.Port = port
	natsSvr := natsserver.RunServer(&opt)
	natsSvr.Start()
	defer natsSvr.Shutdown()

	nc, err := nats.Connect(
		fmt.Sprintf("nats://%s:%d", opt.Host, opt.Port),
	)
	assert.NoError(t, err)
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	assert.NoError(t, err)
	defer ec.Close()

	pub, err := NewPublisher(ec)
	assert.NoError(t, err)

	err = pub.UpdateUser(&User{
		ID: "sample",
	})
	assert.NoError(t, err)
}