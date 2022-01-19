package sqlstore_test

import (
	"testing"

	"paulTabaco/http-rest-api/internal/app/model"
	"paulTabaco/http-rest-api/internal/app/store"
	"paulTabaco/http-rest-api/internal/app/store/sqlstore"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)

	// Testify plugin
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	email := "nonExistinginDB@mail.org"

	_, err := s.User().FindByEmail(email)
	// Testify plugin
	// Case of serching for unexisting user should be negative - ok
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email
	// Case of serching for existing user should be positive - ok
	s.User().Create(u)
	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
