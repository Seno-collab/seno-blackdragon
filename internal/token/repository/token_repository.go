package repository

import (
	"context"

	"seno-blackdragon/internal/token/model"
)

type TokenRepository interface {
	Create(ctx context.Context, t *model.Token) error
	FindByID(ctx context.Context, id int64) (*model.Token, error)
}
