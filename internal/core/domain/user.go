package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/pkg/validator"
)

type User struct {
	ID          uuid.UUID
	Version     int
	Name        string
	Surname     string
	PhoneNumber string
	CreatedAt   time.Time
}

func (u User) Validate() error {
	var v Validation

	v.Add("name", validator.ValidateLen(u.Name, NameMinLen, NameMaxLen))
	v.Add("surname", validator.ValidateLen(u.Surname, SurnameMinLen, SurnameMaxLen))
	v.Add("phone_number", validator.ValidatePhoneNumber(u.PhoneNumber))

	return v.Err()
}

type UserNullable struct {
	Name        Nullable[string]
	Surname     Nullable[string]
	PhoneNumber Nullable[string]
}

func (u UserNullable) Validate() error {
	var v Validation

	if !u.Name.Set && !u.Surname.Set && !u.PhoneNumber.Set {
		v.Addf("", "patch must contain at least one field")
		return v.Err()
	}

	v.AddNullableLen("name", u.Name, NameMinLen, NameMaxLen)
	v.AddNullableLen("surname", u.Surname, SurnameMinLen, SurnameMaxLen)
	v.AddNullableFunc("phone_number", u.PhoneNumber, validator.ValidatePhoneNumber)

	return v.Err()
}

func (u User) WithPatch(patch *UserNullable) (*User, bool) {
	changed := applyNotNull(&u.Name, patch.Name)
	changed = applyNotNull(&u.Surname, patch.Surname) || changed
	changed = applyNotNull(&u.PhoneNumber, patch.PhoneNumber) || changed

	return &u, changed
}
