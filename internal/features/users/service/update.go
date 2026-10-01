package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s UsersService) PatchUser(ctx context.Context, userID uuid.UUID, patch *domain.UserNullable) (int, error) {
	if err := patch.Validate(); err != nil {
		return 0, domain.NewBadRequestError("invalid user data", err)
	}

	origUser, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("get user: %w", err)
	}

	patchedUser, err := origUser.WithPatch(patch)
	if err != nil {
		return 0, domain.NewBadRequestError("nothing to update", err)
	}

	if err := s.usersRepository.UpdateUser(ctx, userID, patchedUser); err != nil {
		return 0, fmt.Errorf("update user: %w", err)
	}

	return patchedUser.Version + 1, nil
}
