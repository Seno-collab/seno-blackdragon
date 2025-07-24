package repository

import (
	"context"
	"errors"
	"sync"

	"seno-blackdragon/internal/user/model"
)

type InMemory struct {
	mu   sync.Mutex
	data map[int64]*model.User
	next int64
}

func NewInMemory() *InMemory {
	return &InMemory{data: make(map[int64]*model.User)}
}

func (r *InMemory) Create(ctx context.Context, u *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	u.ID = r.next
	r.data[u.ID] = u
	return nil
}

func (r *InMemory) FindByID(ctx context.Context, id int64) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.data[id]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}
