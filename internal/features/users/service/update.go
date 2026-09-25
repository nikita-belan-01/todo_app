package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s UsersService) PatchUser(ctx context.Context, userID uuid.UUID, user *domain.UserNullable) error {
	origUser, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	origUser.ApplyPatch(user)

	if err := s.usersRepository.UpdateUser(ctx, userID, origUser); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}
