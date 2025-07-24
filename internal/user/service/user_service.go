package service

import (
	"context"

	"seno-blackdragon/internal/user/model"
	"seno-blackdragon/internal/user/repository"
)

type UserService interface {
	Register(ctx context.Context, name string) (*model.User, error)
	Get(ctx context.Context, id int64) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func New(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, name string) (*model.User, error) {
	user := &model.User{Name: name}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Get(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}
