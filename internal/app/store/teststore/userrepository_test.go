package teststore_test

import (
	"paulTabaco/http-rest-api/internal/app/model"
	"paulTabaco/http-rest-api/internal/app/store"
	"testing"

	"paulTabaco/http-rest-api/internal/app/store/teststore"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {

	s := teststore.New()
	u := model.TestUser1(t)

	// Testify plugin
	assert.NoError(t, s.UserRep().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindById(t *testing.T) {
	s := teststore.New()
	u1 := model.TestUser1(t)
	s.UserRep().Create(u1)
	u2, err := s.UserRep().FindById(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	email := "nonExistinginDB@mail.org"
	_, err := s.UserRep().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser1(t)
	u.Email = email
	s.UserRep().Create(u)
	u, err = s.UserRep().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
