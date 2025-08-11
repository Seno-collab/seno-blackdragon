package repository

import (
	"context"
	"seno-blackdragon/internal/db/user"
	"seno-blackdragon/pkg/utils"
)

type AuthRepo struct {
	q *user.Queries
}

func NewAuthRepo(db user.DBTX) *AuthRepo {
	q := user.New(db)
	return &AuthRepo{q: q}
}

func (r *AuthRepo) GetUserByMail(ctx context.Context, email string) (user.User, error) {
	return r.q.GetUserByEmail(ctx, utils.StringToPgTypeText(email))
}
