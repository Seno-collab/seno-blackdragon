package service

import (
	"context"
	"seno-blackdragon/internal/repository"

	"go.uber.org/zap"
)

type AuthService struct {
	authRepo *repository.AuthRepo
	logger   *zap.Logger
}

func NewAuthService(authRepo *repository.AuthRepo, logger *zap.Logger) *AuthService {
	return &AuthService{authRepo: authRepo, logger: logger}
}

func (as *AuthService) Login(ctx context.Context, email string, password string) (string, string, error) {
	_, err := as.authRepo.GetUserByMail(ctx, email)
	if err != nil {
		as.logger.Error("User not found", zap.String("email", email), zap.Error(err))
		return "", "", err
	}
	return "hello", "heloo", nil
}
