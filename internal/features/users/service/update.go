package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s UsersService) PatchUser(ctx context.Context, userID uuid.UUID, patch *domain.UserNullable) (int, error) {
	if err := patch.Validate(); err != nil {
		return 0, fmt.Errorf("validate: %w", err)
	}

	origUser, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("get: %w", err)
	}

	if origUser.Version != patch.Version {
		return 0, domain.NewAppError(domain.KindConflict, "user was modified, re-read it", domain.ErrVersionConflict)
	}

	patchedUser, changed := origUser.WithPatch(patch)
	if !changed {
		return origUser.Version, nil
	}

	newVersion, err := s.usersRepository.UpdateUser(ctx, patchedUser)
	if err != nil {
		return 0, fmt.Errorf("update: %w", err)
	}

	return newVersion, nil
}
