package user

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("user not found")

type MemoryRepository struct {
	users []User
	mu    sync.RWMutex
}

func NewMemoryRepository(users []User) *MemoryRepository {
	return &MemoryRepository{
		users: append([]User(nil), users...),
	}
}

func (m *MemoryRepository) List(_ context.Context) ([]User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return append([]User(nil), m.users...), nil
}

func (m *MemoryRepository) FindByID(_ context.Context, id string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, candidate := range m.users {
		if candidate.ID == id {
			return &candidate, nil
		}
	}

	return &User{}, ErrNotFound
}

var _ Repository = (*MemoryRepository)(nil)
