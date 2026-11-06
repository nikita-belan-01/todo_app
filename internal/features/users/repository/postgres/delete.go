package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r UsersRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	cmd, err := r.pool.Exec(
		ctx,
		`DELETE FROM todo_app.users WHERE id = $1;`,
		userID)
	if err != nil {
		return fmt.Errorf("exec: %w", utils.TranslateError(err))
	}

	if cmd.RowsAffected() == 0 {
		return domain.NewNotFoundError(domain.ErrUserNotFound)
	}

	return nil
}
