package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r TasksRepository) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	cmd, err := r.pool.Exec(
		ctx,
		`DELETE FROM todo_app.tasks WHERE id = $1`,
		taskID)
	if err != nil {
		return fmt.Errorf("exec: %w", utils.TranslateError(err))
	}

	if cmd.RowsAffected() == 0 {
		return domain.NewAppError(domain.KindNotFound, "not found", domain.ErrTaskNotFound)
	}

	return nil
}
