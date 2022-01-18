package model

import "testing"

/// Make Fake User. To evoide making new faked user each time at begin a new test

// TestUser ...
func TestUser(t *testing.T) *User {
	return &User{
		Email:    "testUser1@mail.org",
		Password: "fakedPassword",
	}
}
