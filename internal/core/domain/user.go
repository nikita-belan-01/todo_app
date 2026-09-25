package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Version     int
	Name        string
	Surname     string
	PhoneNumber string
}

type UserNullable struct {
	Name        Nullable[string]
	Surname     Nullable[string]
	PhoneNumber Nullable[string]
}

func (u *User) ApplyPatch(user *UserNullable) {
	if user.Name.Set && user.Name.Value != nil {
		u.Name = *user.Name.Value
	}

	if user.Surname.Set && user.Surname.Value != nil {
		u.Surname = *user.Surname.Value
	}

	if user.PhoneNumber.Set && user.PhoneNumber.Value != nil {
		u.PhoneNumber = *user.PhoneNumber.Value
	}
}
