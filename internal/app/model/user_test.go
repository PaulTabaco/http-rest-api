package model_test

import (
	"paulTabaco/http-rest-api/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
