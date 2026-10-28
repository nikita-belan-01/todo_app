package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r TasksRepository) CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {
	var id uuid.UUID
	if err := r.pool.QueryRow(
		ctx,
		`INSERT INTO todo_app.tasks (
									user_id,
									title,
									description)
		VALUES ($1, $2, $3) RETURNING id`,
		task.UserID,
		task.Title,
		task.Description,
	).Scan(&id); err != nil {
		return uuid.Nil, fmt.Errorf("scan row: %w", utils.TranslateError(err))
	}

	return id, nil
}
