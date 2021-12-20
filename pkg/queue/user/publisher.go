package user

import (
	"github.com/nats-io/nats.go"
	"time"
)

type publisher struct {
	ready bool
	ec    *nats.EncodedConn
}

type Publisher interface {
	IsReady() bool
	UpdateUser(u *User) error
}

func NewPublisher(ec *nats.EncodedConn) (Publisher, error) {
	pub := publisher{
		ready: true,
		ec:    ec,
	}
	go func() {
		tic := time.Tick(time.Minute * 1)
		for range tic {
			if err := ec.Flush(); err != nil {
				pub.ready = false
			}
		}
	}()
	return &pub, nil
}

func (s *publisher) IsReady() bool {
	return s.ready
}

func (s *publisher) UpdateUser(u *User) error {
	return s.ec.Publish(UpdateUserSubject, u)
}
