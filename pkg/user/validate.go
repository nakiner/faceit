package user

import (
	"github.com/pkg/errors"
)

type validator interface {
	Validate() error
}

func validate(req interface{}) error {
	if val, ok := interface{}(req).(validator); ok {
		if err := val.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (r CreateUserRequest) Validate() error {
	if r.Password != r.PasswordConfirm {
		return errors.Wrap(ErrBadRequest, "passwords does not match")
	}
	return nil
}

func (r DeleteUserRequest) Validate() error {
	if len(r.Id) < 1 {
		return errors.Wrap(ErrBadRequest, "id cannot be empty")
	}
	return nil
}

func (r User) Validate() error {
	if len(r.Id) < 1 {
		return errors.Wrap(ErrBadRequest, "id cannot be empty")
	}
	return nil
}

func (r GetUsersRequest) Validate() error {
	if r.Limit > 500 {
		return errors.Wrap(ErrBadRequest, "limit should not be greater then 500")
	}
	return nil
}
