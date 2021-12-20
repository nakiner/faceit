//go:build integration && !unit
// +build integration,!unit

package integration

import (
	"fmt"
	"github.com/nakiner/faceit/configs"
	"github.com/nakiner/faceit/pkg/queue/user"
	"log"
	"testing"
	"time"

	natsCl "github.com/nakiner/faceit/pkg/store/nats"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNatsPublisherUserServiceUpdateUser(t *testing.T) {
	cfg := configs.NewConfig()
	err := cfg.Read()
	if err != nil {
		log.Fatal(err)
	}

	nc, err := nats.Connect(
		fmt.Sprintf("nats://%s:%d", cfg.Nats.Host, cfg.Nats.Port),
		nats.ReconnectWait(time.Millisecond*time.Duration(cfg.Nats.WaitLimit)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	p, err := user.NewPublisher(ec)
	if err != nil {
		log.Fatal(err)
	}
	err = p.UpdateUser(&user.User{Nickname: "sample"})
	require.NoError(t, err)
}

func TestNatsSubscriberUserServiceUpdateUser(t *testing.T) {
	cfg := configs.NewConfig()
	err := cfg.Read()
	if err != nil {
		log.Fatal(err)
	}

	nc, err := natsCl.NewClient(&cfg.Nats)
	assert.NoError(t, err)
	defer nc.Close()

	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ch := make(chan string)
	s := user.NewSubscriber(nc)
	err = s.UpdateUser(func(u *user.User) {
		ch <- u.ID
	})
	require.NoError(t, err)

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	us := user.User{
		ID: "sample",
	}

	p, err := user.NewPublisher(ec)
	assert.NoError(t, err)

	for i := 0; i < 5; i++ {
		err = p.UpdateUser(&us)
		assert.NoError(t, err)
	}

	time.Sleep(time.Second * 1)

	res := <-ch

	assert.Equal(t, us.ID, res)
}
