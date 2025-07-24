package service

import (
	"context"

	"seno-blackdragon/internal/dragon/model"
	"seno-blackdragon/internal/dragon/repository"
)

type DragonService interface {
	Create(ctx context.Context, ownerID int64, name string) (*model.Dragon, error)
	Get(ctx context.Context, id int64) (*model.Dragon, error)
}

type dragonService struct {
	repo repository.DragonRepository
}

func New(repo repository.DragonRepository) DragonService {
	return &dragonService{repo: repo}
}

func (s *dragonService) Create(ctx context.Context, ownerID int64, name string) (*model.Dragon, error) {
	d := &model.Dragon{OwnerID: ownerID, Name: name}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *dragonService) Get(ctx context.Context, id int64) (*model.Dragon, error) {
	return s.repo.FindByID(ctx, id)
}
