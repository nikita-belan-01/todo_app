package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s UsersService) PatchUser(ctx context.Context, userID uuid.UUID, patch *domain.UserNullable) (int, error) {
	if err := patch.Validate(); err != nil {
		return 0, domain.NewAppError(domain.KindInvalidArgument, "invalid data", err)
	}

	origUser, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("get: %w", err)
	}

	patchedUser, err := origUser.WithPatch(patch)
	if err != nil {
		return 0, fmt.Errorf("patching: %w", err)
	}

	newVersion, err := s.usersRepository.UpdateUser(ctx, patchedUser)
	if err != nil {
		return 0, fmt.Errorf("update: %w", err)
	}

	return newVersion, nil
}
