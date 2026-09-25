package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (r UsersRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	cmd, err := r.pool.Exec(
		ctx,
		`DELETE FROM todo_app.users WHERE id = $1;`,
		userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return domain.NewNotFoundError("user not found", domain.ErrUserNotFound)
	}

	return nil
}
