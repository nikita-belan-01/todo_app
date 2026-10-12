package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s UsersService) GetUsers(ctx context.Context, page, limit int) ([]domain.User, error) {
	data, err := s.usersRepository.GetUsers(ctx, page, limit)
	if err != nil {
		return nil, err
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
