package service

import (
	"context"

	"seno-blackdragon/internal/token/model"
	"seno-blackdragon/internal/token/repository"
)

type TokenService interface {
	Create(ctx context.Context, userID int64, value string) (*model.Token, error)
	Get(ctx context.Context, id int64) (*model.Token, error)
}

type tokenService struct {
	repo repository.TokenRepository
}

func New(repo repository.TokenRepository) TokenService {
	return &tokenService{repo: repo}
}

func (s *tokenService) Create(ctx context.Context, userID int64, value string) (*model.Token, error) {
	t := &model.Token{UserID: userID, Value: value}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *tokenService) Get(ctx context.Context, id int64) (*model.Token, error) {
	return s.repo.FindByID(ctx, id)
}
