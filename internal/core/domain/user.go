package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/pkg/validator"
)

type User struct {
	ID          uuid.UUID
	Version     int
	Name        string
	Surname     string
	PhoneNumber string
}

func (u User) Validate() error {
	var errs []error
	if err := validator.ValidateLen(u.Name, 3, 100); err != nil {
		errs = append(errs, fmt.Errorf("name: %w", err))
	}

	if err := validator.ValidateLen(u.Surname, 3, 100); err != nil {
		errs = append(errs, fmt.Errorf("surname: %w", err))
	}

	if err := validator.ValidatePhoneNumber(u.PhoneNumber); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type UserNullable struct {
	Name        Nullable[string]
	Surname     Nullable[string]
	PhoneNumber Nullable[string]
}

func (u UserNullable) Validate() error {
	var errs []error
	if u.Name.Set {
		if u.Name.Value == nil {
			errs = append(errs, fmt.Errorf("name can't be patched to NULL: %w", ErrInvalidArgument))
		}
		if u.Name.Value != nil {
			if err := validator.ValidateLen(*u.Name.Value, 3, 100); err != nil {
				errs = append(errs, fmt.Errorf("name: %w", err))
			}
		}
	}

	if u.Surname.Set == true {
		if u.Surname.Value == nil {
			errs = append(errs, fmt.Errorf("surname can't be patched to NULL: %w", ErrInvalidArgument))
		}
		if u.Surname.Value != nil {
			if err := validator.ValidateLen(*u.Surname.Value, 3, 100); err != nil {
				errs = append(errs, fmt.Errorf("surname: %w", err))
			}
		}
	}

	if u.PhoneNumber.Set == true {
		if u.PhoneNumber.Value == nil {
			errs = append(errs, fmt.Errorf("phone number can't be patched to NULL: %w", ErrInvalidArgument))
		}
		if u.PhoneNumber.Value != nil {
			if err := validator.ValidatePhoneNumber(*u.PhoneNumber.Value); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
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
