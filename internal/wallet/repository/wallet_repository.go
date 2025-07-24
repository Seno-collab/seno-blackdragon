package repository

import (
	"context"

	"seno-blackdragon/internal/wallet/model"
)

type WalletRepository interface {
	Create(ctx context.Context, w *model.Wallet) error
	FindByID(ctx context.Context, id int64) (*model.Wallet, error)
}
