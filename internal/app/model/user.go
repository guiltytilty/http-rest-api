package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	is "github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       *int      `json:"age,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Age   *int    `json:"age"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&u.Age, validation.Min(0), validation.Max(150)),
	)
}
