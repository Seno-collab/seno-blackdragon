package service

import (
	"context"

	"seno-blackdragon/internal/wallet/model"
	"seno-blackdragon/internal/wallet/repository"
)

type WalletService interface {
	Create(ctx context.Context, userID int64) (*model.Wallet, error)
	Get(ctx context.Context, id int64) (*model.Wallet, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func New(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) Create(ctx context.Context, userID int64) (*model.Wallet, error) {
	w := &model.Wallet{UserID: userID}
	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *walletService) Get(ctx context.Context, id int64) (*model.Wallet, error) {
	return s.repo.FindByID(ctx, id)
}
