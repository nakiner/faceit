package user

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type subscriber struct {
	ready bool
	nc    *nats.Conn
}

type UpdateUserHandler func(u *User)

type Subscriber interface {
	UpdateUser(fn UpdateUserHandler) error
}

func NewSubscriber(nc *nats.Conn) Subscriber {
	return &subscriber{
		ready: true,
		nc:    nc,
	}
}

func (s *subscriber) UpdateUser(fn UpdateUserHandler) error {
	if _, err := s.nc.QueueSubscribe(UpdateUserSubject, Queue, func(msg *nats.Msg) {
		dec, err := decodeNATSUserRequest(context.Background(), msg)
		if err == nil {
			user := dec.(User)
			fn(&user)
		}
	}); err != nil {
		return err
	}

	return nil
}

func decodeNATSUserRequest(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var request User
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		return nil, err
	}
	return request, nil
}
