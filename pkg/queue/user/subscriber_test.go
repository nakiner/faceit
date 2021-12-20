package user

import (
	natsCl "github.com/nakiner/faceit/pkg/store/nats"
	natsserver "github.com/nats-io/nats-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSubscriber(t *testing.T) {
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

	sub := NewSubscriber(nc)
	assert.NotNil(t, sub)
}

func TestSubscriber_UpdateUser(t *testing.T) {
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

	sub := NewSubscriber(nc)
	assert.NotNil(t, sub)

	err = sub.UpdateUser(func(u *User) {})
	assert.NoError(t, err)
}
