package sqlstore

import (
	"database/sql"
	"paulTabaco/http-rest-api/internal/app/store"

	_ "github.com/lib/pq" // ...
)

// Store ...
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// UserRep ...
// for outer using user repository - store.UserRep().Create() , .FindByEmail() ...
func (s *Store) UserRep() store.UserRepositoryInterface {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
