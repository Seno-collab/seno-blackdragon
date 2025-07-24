package repository

import (
	"context"
	"errors"
	"sync"

	"seno-blackdragon/internal/token/model"
)

type InMemory struct {
	mu   sync.Mutex
	data map[int64]*model.Token
	next int64
}

func NewInMemory() *InMemory {
	return &InMemory{data: make(map[int64]*model.Token)}
}

func (r *InMemory) Create(ctx context.Context, t *model.Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	t.ID = r.next
	r.data[t.ID] = t
	return nil
}

func (r *InMemory) FindByID(ctx context.Context, id int64) (*model.Token, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.data[id]; ok {
		return t, nil
	}
	return nil, errors.New("not found")
}
