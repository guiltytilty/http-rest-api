package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Name:  "name",
		Email: "name@example.com",
	}
}
