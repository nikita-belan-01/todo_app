package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s UsersService) CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	if err := user.Validate(); err != nil {
		return uuid.Nil, fmt.Errorf("validate: %w", err)
	}

	id, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
