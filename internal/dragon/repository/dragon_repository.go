package repository

import (
	"context"

	"seno-blackdragon/internal/dragon/model"
)

type DragonRepository interface {
	Create(ctx context.Context, d *model.Dragon) error
	FindByID(ctx context.Context, id int64) (*model.Dragon, error)
}
