package domain

import (
	"errors"
	"fmt"
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
	if !u.Name.Set && !u.Surname.Set && !u.PhoneNumber.Set {
		return fmt.Errorf("patch must contain at least one field: %w", ErrInvalidArgument)
	}

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

	if u.Surname.Set {
		if u.Surname.Value == nil {
			errs = append(errs, fmt.Errorf("surname can't be patched to NULL: %w", ErrInvalidArgument))
		}
		if u.Surname.Value != nil {
			if err := validator.ValidateLen(*u.Surname.Value, 3, 100); err != nil {
				errs = append(errs, fmt.Errorf("surname: %w", err))
			}
		}
	}

	if u.PhoneNumber.Set {
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

func (u User) WithPatch(patch *UserNullable) (*User, error) {
	var changed bool

	if patch.Name.Set && patch.Name.Value != nil && *patch.Name.Value != u.Name {
		u.Name = *patch.Name.Value
		changed = true
	}

	if patch.Surname.Set && patch.Surname.Value != nil && *patch.Surname.Value != u.Surname {
		u.Surname = *patch.Surname.Value
		changed = true
	}

	if patch.PhoneNumber.Set && patch.PhoneNumber.Value != nil && *patch.PhoneNumber.Value != u.PhoneNumber {
		u.PhoneNumber = *patch.PhoneNumber.Value
		changed = true
	}

	if !changed {
		return nil, ErrNothingToUpdate
	}

	return &u, nil
}
