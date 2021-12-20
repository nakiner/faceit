//go:build integration && !unit
// +build integration,!unit

package integration

import (
	"github.com/nakiner/faceit/pkg/queue/user"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestNatsPublisherUserServiceUpdateUser(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	p, err := user.NewPublisher(ec)
	if err != nil {
		log.Fatal(err)
	}
	err = p.UpdateUser(&user.User{})
	require.NoError(t, err)
}

func TestNatsSubscriberUserServiceUpdateUser(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	p, err := user.NewPublisher(ec)
	if err != nil {
		log.Fatal(err)
	}
	n := "sample"
	err = p.UpdateUser(&user.User{Nickname: "smpl"})
	if err != nil {
		log.Fatal(err)
	}
	require.NoError(t, err)

	var wg sync.WaitGroup

	wg.Add(1)

	s := user.NewSubscriber(nc)
	err = s.UpdateUser(func(u *user.User) {
		assert.Equal(t, n, u.Nickname)
	})

	wg.Wait()

	require.NoError(t, err)
}
