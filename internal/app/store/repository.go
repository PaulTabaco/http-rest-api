package store

import "paulTabaco/http-rest-api/internal/app/model"

// UserRepositoryInterface ...
type UserRepositoryInterface interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
