package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/config"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/server"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

type UsersHTTPHandler struct {
	Config       *config.Handler
	UsersService UsersService
}

type UsersService interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	PatchUser(ctx context.Context, userID uuid.UUID, user *domain.UserNullable) error
}

func NewUsersHTTPHandler(config *config.Handler, usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		Config:       config,
		UsersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []server.Route {
	const usersPath string = "/users"
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    usersPath,
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    usersPath,
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    fmt.Sprintf("%s/{userId}", usersPath),
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    fmt.Sprintf("%s/{userId}", usersPath),
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    fmt.Sprintf("%s/{userId}", usersPath),
			Handler: h.PatchUser,
		}}
}
