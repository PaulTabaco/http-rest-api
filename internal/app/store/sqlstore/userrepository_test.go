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
	u := model.TestUser1(t)

	// Testify plugin
	assert.NoError(t, s.UserRep().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindById(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u1 := model.TestUser1(t)
	s.UserRep().Create(u1)
	u2, err := s.UserRep().FindById(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u1 := model.TestUser1(t)
	_, err := s.UserRep().FindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.UserRep().Create(u1)
	u2, err := s.UserRep().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
