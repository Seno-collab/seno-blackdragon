package repository

import (
	"context"
	"errors"
	"sync"

	"seno-blackdragon/internal/skill/model"
)

type InMemory struct {
	mu   sync.Mutex
	data map[int64]*model.Skill
	next int64
}

func NewInMemory() *InMemory {
	return &InMemory{data: make(map[int64]*model.Skill)}
}

func (r *InMemory) Create(ctx context.Context, s *model.Skill) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	s.ID = r.next
	r.data[s.ID] = s
	return nil
}

func (r *InMemory) FindByID(ctx context.Context, id int64) (*model.Skill, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.data[id]; ok {
		return s, nil
	}
	return nil, errors.New("not found")
}
