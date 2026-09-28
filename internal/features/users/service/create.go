package service

import (
	"context"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s UsersService) CreateUser(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return domain.NewBadRequestError("validation user", err)
	}

	if err := s.usersRepository.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}
