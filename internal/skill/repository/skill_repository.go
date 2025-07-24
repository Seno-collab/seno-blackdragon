package repository

import (
	"context"

	"seno-blackdragon/internal/skill/model"
)

type SkillRepository interface {
	Create(ctx context.Context, s *model.Skill) error
	FindByID(ctx context.Context, id int64) (*model.Skill, error)
}
