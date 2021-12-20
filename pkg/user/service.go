package user

import (
	"context"
	"github.com/go-kit/log/level"
	userRepository "github.com/nakiner/faceit/internal/repository/user"
	userQueue "github.com/nakiner/faceit/pkg/queue/user"
	"github.com/nakiner/faceit/tools/logging"
	"github.com/pkg/errors"
)

type userService struct {
	repo      userRepository.Repository
	ncUserPub userQueue.Publisher
}

func NewUserService(repo userRepository.Repository, ncUserPub userQueue.Publisher) Service {
	return &userService{
		repo:      repo,
		ncUserPub: ncUserPub,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (resp *CreateUserResponse, err error) {
	id, err := s.repo.Create(ctx, &userRepository.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	})
	if err != nil {
		return nil, errors.Wrap(err, "userService create user err")
	}
	resp = &CreateUserResponse{
		Id: id,
	}
	return resp, nil
}

func (s *userService) GetUsers(ctx context.Context, req *GetUsersRequest) (resp *GetUsersResponse, err error) {
	conditions := userRepository.Conditions{}

	if len(req.Id) > 0 {
		conditions["id"] = req.Id
	}
	if len(req.Country) > 0 {
		conditions["country"] = req.Country
	}
	if len(req.Nickname) > 0 {
		conditions["nickname"] = req.Nickname
	}
	if len(req.FirstName) > 0 {
		conditions["first_name"] = req.FirstName
	}
	if len(req.LastName) > 0 {
		conditions["last_name"] = req.LastName
	}
	if req.Limit < 1 {
		req.Limit = 50
	}
	if req.Offset < 1 {
		req.Offset = 1
	}

	users, err := s.repo.Get(ctx, conditions, req.Limit, req.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "userService GetUsers err")
	}

	data := GetUsersResponse{}

	for _, user := range users {
		data = append(data, User{
			Id:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Password:  user.Password,
			Email:     user.Email,
			Country:   user.Country,
			CreatedAt: user.TimeToString(user.CreatedAt),
			UpdatedAt: user.TimeToString(user.UpdatedAt),
		})
	}

	return &data, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *User) (resp *Status, err error) {
	err = s.repo.Update(ctx, &userRepository.User{
		ID:        req.Id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	})
	if errors.Is(err, userRepository.ErrRowsAffectedEmpty) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "userService update user err")
	}
	go func() {
		if err = s.ncUserPub.UpdateUser(&userQueue.User{
			ID:        req.Id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Nickname:  req.Nickname,
			Password:  req.Password,
			Email:     req.Email,
			Country:   req.Country,
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
		}); err != nil {
			lg := logging.FromContext(context.Background())
			level.Error(lg).Log("msg", "could not pub to channel", "err", err)
		}
	}()
	return &Status{
		Status:  true,
		Message: "OK",
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *DeleteUserRequest) (resp *Status, err error) {
	err = s.repo.Delete(ctx, req.Id)
	if errors.Is(err, userRepository.ErrRowsAffectedEmpty) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "userService deleteById err")
	}
	return &Status{
		Status:  true,
		Message: "OK",
	}, nil
}
