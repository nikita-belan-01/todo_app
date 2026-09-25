package service

import (
	"context"

	"github.com/google/uuid"
)

func (s UsersService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.usersRepository.DeleteUser(ctx, userID); err != nil {
		return err
	}

	return nil
}
