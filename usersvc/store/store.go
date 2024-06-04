package store

import (
	"context"
	"errors"
	"sync"

	"github.com/ficontini/user-search/types"
)

type Store interface {
	Add(context.Context, *types.User) error
	GetByID(context.Context, int64) (*types.User, error)
	GetByIDs(context.Context, []int64) ([]*types.User, error)
	SearchByCriteria(context.Context, types.SearchCriteria) ([]*types.User, error)
}

type InMemoryStore struct {
	mu    sync.RWMutex
	users map[int64]*types.User
}

func NewInMemoryStore() Store {
	store := &InMemoryStore{
		users: make(map[int64]*types.User),
	}

	store.seedData()

	return store
}
func (s *InMemoryStore) Add(_ context.Context, user *types.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryStore) GetByID(ctx context.Context, id int64) (*types.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user := s.users[id]
	if user != nil {
		return user, nil
	}
	return nil, errors.New("user not found")
}
func (s *InMemoryStore) GetByIDs(ctx context.Context, ids []int64) ([]*types.User, error) {
	var users []*types.User
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		user := s.users[id]
		if user != nil {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		return nil, errors.New("users not found")
	}
	return users, nil
}

func (s *InMemoryStore) seedData() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range GenerateRandomUsers() {
		s.users[user.ID] = user
	}
}
func (s *InMemoryStore) SearchByCriteria(_ context.Context, criteria types.SearchCriteria) ([]*types.User, error) {
	var users []*types.User
	for _, u := range s.users {
		if criteria.Meets(u) {
			users = append(users, u)
		}
	}
	return users, nil
}
