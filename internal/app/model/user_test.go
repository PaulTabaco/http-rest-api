package model_test

import (
	"paulTabaco/http-rest-api/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid case",
			u: func() *model.User {
				return model.TestUser1(t)
			},
			isValid: true,
		},
		{
			name: "with encrypted password case",
			u: func() *model.User {
				u := model.TestUser1(t)
				u.Password = ""
				u.EncryptedPassword = "nonEmptyEncryptedPassword"
				return u
			},
			isValid: true,
		},
		{
			name: "empty email case",
			u: func() *model.User {
				u := model.TestUser1(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email case",
			u: func() *model.User {
				u := model.TestUser1(t)
				u.Email = "inv@"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *model.User {
				u := model.TestUser1(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short email case",
			u: func() *model.User {
				u := model.TestUser1(t)
				u.Password = "1234"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}

	u := model.TestUser1(t)
	assert.NoError(t, u.Validate())

}

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser1(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
