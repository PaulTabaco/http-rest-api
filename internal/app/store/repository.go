package store

import "paulTabaco/http-rest-api/internal/app/model"

// UserRepository
type UserRepository interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
