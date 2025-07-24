package repository

import (
	"context"

	"seno-blackdragon/internal/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	FindByID(ctx context.Context, id int64) (*model.User, error)
}
