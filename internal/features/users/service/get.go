package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s UsersService) GetUsers(ctx context.Context, pagination *domain.Pagination) ([]domain.User, error) {
	if err := pagination.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	data, err := s.usersRepository.GetUsers(ctx, pagination)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return data, nil
}

func (s UsersService) GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
