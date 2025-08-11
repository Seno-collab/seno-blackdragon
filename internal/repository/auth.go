package repository

import (
	"context"
	"database/sql"
	"errors"
	"seno-blackdragon/internal/db/user"
	"seno-blackdragon/pkg/response_status"
	"seno-blackdragon/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepo struct {
	q *user.Queries
}

type UserModel struct {
	ID           uuid.UUID
	FullName     string
	Bio          string
	Email        string
	PasswordHash string
}

func NewAuthRepo(db user.DBTX) *AuthRepo {
	q := user.New(db)
	return &AuthRepo{q: q}
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	row, err := r.q.GetUserByEmail(ctx, utils.StringToPgTypeText(email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response_status.ErrUserNotFound
		}
		return nil, err
	}
	return &row, nil
}

func (r *AuthRepo) RegisterUser(ctx context.Context, email, fullName, bio, password_hash string) (pgtype.UUID, error) {
	user := user.AddUserParams{
		FullName:     fullName,
		Email:        utils.StringToPgTypeText(email),
		Bio:          utils.StringToPgTypeText(bio),
		PasswordHash: utils.StringToPgTypeText(password_hash),
	}
	id, error := r.q.AddUser(ctx, user)
	return id, error
}

func (r AuthRepo) CreateUser(ctx context.Context, u *UserModel) (uuid.UUID, error) {
	params := user.AddUserParams{
		FullName:     u.FullName,
		Email:        utils.StringToPgTypeText(u.Email),
		Bio:          utils.StringToPgTypeText(u.Bio),
		PasswordHash: utils.StringToPgTypeText(u.PasswordHash),
	}

	id, err := r.q.AddUser(ctx, params)
	if err != nil {
		return uuid.Nil, err
	}
	return utils.UUIDFromPgUUID(id), nil
}
