package http

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/types"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/pkg/validator"
)

const userIDPathKey = "userId"

type createUserRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number"`
}

func (r createUserRequest) Validate() error {
	var errs []error
	if err := validator.ValidateLen(r.Name, 3, 100); err != nil {
		errs = append(errs, fmt.Errorf("name: %w", err))
	}

	if err := validator.ValidateLen(r.Surname, 3, 100); err != nil {
		errs = append(errs, fmt.Errorf("surname: %w", err))
	}

	if err := validator.ValidatePhoneNumber(r.PhoneNumber); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func BuildDomainUser(name, surname, phoneNumber string) *domain.User {
	return &domain.User{
		Name:        name,
		Surname:     surname,
		PhoneNumber: phoneNumber,
	}
}

type userResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	PhoneNumber string    `json:"phone_number"`
}

func BuildUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		PhoneNumber: user.PhoneNumber,
	}
}

func BuildUsersResponse(users []domain.User) []userResponse {
	response := make([]userResponse, 0, len(users))
	for _, user := range users {
		response = append(response, BuildUserResponse(&user))
	}

	return response
}

type patchUserRequest struct {
	Name        types.Nullable[string] `json:"name"`
	Surname     types.Nullable[string] `json:"surname"`
	PhoneNumber types.Nullable[string] `json:"phone_number"`
}

func (r patchUserRequest) Validate() error {
	var errs []error
	if r.Name.Set {
		if r.Name.Value == nil {
			errs = append(errs, fmt.Errorf("name can't be patched to NULL: %w", domain.ErrInvalidArgument))
		}
		if r.Name.Value != nil {
			if err := validator.ValidateLen(*r.Name.Value, 3, 100); err != nil {
				errs = append(errs, fmt.Errorf("name: %w", err))
			}
		}
	}

	if r.Surname.Set == true {
		if r.Surname.Value == nil {
			errs = append(errs, fmt.Errorf("surname can't be patched to NULL: %w", domain.ErrInvalidArgument))
		}
		if r.Surname.Value != nil {
			if err := validator.ValidateLen(*r.Surname.Value, 3, 100); err != nil {
				errs = append(errs, fmt.Errorf("surname: %w", err))
			}
		}
	}

	if r.PhoneNumber.Set == true {
		if r.PhoneNumber.Value == nil {
			errs = append(errs, fmt.Errorf("phone number can't be patched to NULL: %w", domain.ErrInvalidArgument))
		}
		if r.PhoneNumber.Value != nil {
			if err := validator.ValidatePhoneNumber(*r.PhoneNumber.Value); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

func BuildDomainNullableUser(name, surname, phoneNumber types.Nullable[string]) *domain.UserNullable {
	return &domain.UserNullable{
		Name:        name.Nullable,
		Surname:     surname.Nullable,
		PhoneNumber: phoneNumber.Nullable,
	}
}
