package http

import (
	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/types"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

const userIDPathKey = "userId"

type createUserRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number"`
}

func (r createUserRequest) BuildDomain() domain.User {
	return domain.User{
		Name:        r.Name,
		Surname:     r.Surname,
		PhoneNumber: r.PhoneNumber,
	}
}

type userResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	PhoneNumber string    `json:"phone_number"`
}

type userIDResponse struct {
	ID uuid.UUID `json:"id"`
}

type userPatchResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`
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

func (r patchUserRequest) BuildDomain() domain.UserNullable {
	return domain.UserNullable{
		Name:        r.Name.Nullable,
		Surname:     r.Surname.Nullable,
		PhoneNumber: r.PhoneNumber.Nullable,
	}
}
