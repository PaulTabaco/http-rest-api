package teststore

import (
	"paulTabaco/http-rest-api/internal/app/model"
	"paulTabaco/http-rest-api/internal/app/store"
)

// Store ...
type Store struct {
	userRepository *UserRepository
}

// New ...
func New() *Store {
	return &Store{}
}

// User ...
// for outer using user repository - store.User().Create() , .FindByEmail() ...
func (s *Store) UserRep() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.userRepository
}
