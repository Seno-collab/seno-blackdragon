package repository

import (
	"context"
	"errors"
	"sync"

	"seno-blackdragon/internal/dragon/model"
)

type InMemory struct {
	mu   sync.Mutex
	data map[int64]*model.Dragon
	next int64
}

func NewInMemory() *InMemory {
	return &InMemory{data: make(map[int64]*model.Dragon)}
}

func (r *InMemory) Create(ctx context.Context, d *model.Dragon) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	d.ID = r.next
	r.data[d.ID] = d
	return nil
}

func (r *InMemory) FindByID(ctx context.Context, id int64) (*model.Dragon, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d, ok := r.data[id]; ok {
		return d, nil
	}
	return nil, errors.New("not found")
}
