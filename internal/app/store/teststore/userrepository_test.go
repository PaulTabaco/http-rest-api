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
	u1 := model.TestUser1(t)
	_, err := s.UserRep().FindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.UserRep().Create(u1)
	u2, err := s.UserRep().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
