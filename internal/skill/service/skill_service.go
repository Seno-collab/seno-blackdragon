package service

import (
	"context"

	"seno-blackdragon/internal/skill/model"
	"seno-blackdragon/internal/skill/repository"
)

type SkillService interface {
	Create(ctx context.Context, name, description string) (*model.Skill, error)
	Get(ctx context.Context, id int64) (*model.Skill, error)
}

type skillService struct {
	repo repository.SkillRepository
}

func New(repo repository.SkillRepository) SkillService {
	return &skillService{repo: repo}
}

func (s *skillService) Create(ctx context.Context, name, description string) (*model.Skill, error) {
	sk := &model.Skill{Name: name, Description: description}
	if err := s.repo.Create(ctx, sk); err != nil {
		return nil, err
	}
	return sk, nil
}

func (s *skillService) Get(ctx context.Context, id int64) (*model.Skill, error) {
	return s.repo.FindByID(ctx, id)
}
