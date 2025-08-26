package repository

import (
	"context"
	"errors"
	"seno-blackdragon/internal/db/user"
	"seno-blackdragon/pkg/enum"
	"seno-blackdragon/pkg/utils"

	"github.com/google/uuid"
)

type UserRepo struct {
	q *user.Queries
}

type UserModel struct {
	ID           uuid.UUID
	FullName     string
	Bio          string
	Email        string
	PasswordHash string
}

func NewUserRepo(db user.DBTX) *UserRepo {
	q := user.New(db)
	return &UserRepo{q: q}
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (*UserModel, error) {
	row, err := ur.q.GetUserByEmail(ctx, utils.PgTextFromString(email))
	if err != nil {
		if errors.Is(err, enum.ErrUserNotFound) {
			return nil, enum.ErrUserNotFound
		}
		return nil, err
	}
	user := &UserModel{
		ID:           utils.UUIDFromPgUUID(row.ID),
		FullName:     row.FullName,
		Bio:          utils.StringFromPgText(row.Bio),
		PasswordHash: utils.StringFromPgText(row.PasswordHash),
	}
	return user, nil
}

func (ur UserRepo) CreateUser(ctx context.Context, u *UserModel) (uuid.UUID, error) {
	params := user.AddUserParams{
		FullName:     u.FullName,
		Email:        utils.PgTextFromString(u.Email),
		Bio:          utils.PgTextFromString(u.Bio),
		PasswordHash: utils.PgTextFromString(u.PasswordHash),
	}

	id, err := ur.q.AddUser(ctx, params)
	if err != nil {
		return uuid.Nil, err
	}
	return utils.UUIDFromPgUUID(id), nil
}

func (ur UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*UserModel, error) {
	row, err := ur.q.GetUserByID(ctx, utils.PgUUIDFromUUID(id))
	if err != nil {
		return nil, err
	}
	u := &UserModel{
		ID:           utils.UUIDFromPgUUID(row.ID),
		FullName:     row.FullName,
		Bio:          utils.StringFromPgText(row.Bio),
		PasswordHash: utils.StringFromPgText(row.PasswordHash),
	}
	return u, nil
}
