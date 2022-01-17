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

	u, err := s.User().Create(&model.User{
		Email: "testUser1@mail.com",
	})

	// Testify plugin
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	email := "testUser0000@mail.com"

	_, err := s.User().FindByEmail(email)

	// Testify plugin
	// Case of serching for unexisting user should be error
	assert.Error(t, err)

	// Case of serching for unexisting user should be error
	s.User().Create(&model.User{
		Email: "testUser0000@mail.com",
	})
	u, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
