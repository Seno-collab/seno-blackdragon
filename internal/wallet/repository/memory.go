package repository

import (
	"context"
	"errors"
	"sync"

	"seno-blackdragon/internal/wallet/model"
)

type InMemory struct {
	mu   sync.Mutex
	data map[int64]*model.Wallet
	next int64
}

func NewInMemory() *InMemory {
	return &InMemory{data: make(map[int64]*model.Wallet)}
}

func (r *InMemory) Create(ctx context.Context, w *model.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	w.ID = r.next
	r.data[w.ID] = w
	return nil
}

func (r *InMemory) FindByID(ctx context.Context, id int64) (*model.Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w, ok := r.data[id]; ok {
		return w, nil
	}
	return nil, errors.New("not found")
}
