package model_test

import (
	"http-rest-api/internal/app/model"
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
			name: "valid without age",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "valid with age",
			u: func() *model.User {
				u := model.TestUser(t)
				age := 52
				u.Age = &age
				return u
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "empty name",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Name = ""
				return u
			},
			isValid: false,
		},
		{
			name: "too long name",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Name = "longnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongname"
				return u
			},
			isValid: false,
		},
		{
			name: "too old",
			u: func() *model.User {
				u := model.TestUser(t)
				age := 152
				u.Age = &age
				return u
			},
			isValid: false,
		},
		{
			name: "too young",
			u: func() *model.User {
				u := model.TestUser(t)
				age := -1
				u.Age = &age
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
}
