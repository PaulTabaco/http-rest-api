package store_test

import (
	"paulTabaco/http-rest-api/internal/app/model"
	"paulTabaco/http-rest-api/internal/app/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	u, err := s.User().Create(model.TestUser(t))

	// Testify plugin
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	email := "testUser1@mail.org"

	_, err := s.User().FindByEmail(email)
	// Testify plugin
	// Case of serching for unexisting user should be negative - ok
	assert.Error(t, err)

	u := model.TestUser(t)
	u.Email = email
	// Case of serching for existing user should be positive - ok
	s.User().Create(u)
	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
