package user

import (
	natsCl "github.com/nakiner/faceit/pkg/store/nats"
	natsserver "github.com/nats-io/nats-server/test"
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

	nc, err := natsCl.NewClient(&natsCl.Config{
		Host: opt.Host,
		Port: opt.Port,
	})
	assert.NoError(t, err)
	defer nc.Close()

	ec, err := natsCl.NewEncodedClient(nc)
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

	nc, err := natsCl.NewClient(&natsCl.Config{
		Host: opt.Host,
		Port: opt.Port,
	})
	assert.NoError(t, err)
	defer nc.Close()

	ec, err := natsCl.NewEncodedClient(nc)
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

	nc, err := natsCl.NewClient(&natsCl.Config{
		Host: opt.Host,
		Port: opt.Port,
	})
	assert.NoError(t, err)
	defer nc.Close()

	ec, err := natsCl.NewEncodedClient(nc)
	assert.NoError(t, err)
	defer ec.Close()

	pub, err := NewPublisher(ec)
	assert.NoError(t, err)

	err = pub.UpdateUser(&User{
		ID: "sample",
	})
	assert.NoError(t, err)
}
