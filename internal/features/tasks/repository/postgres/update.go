package postgres

import (
	"context"
	"fmt"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r TasksRepository) UpdateTask(ctx context.Context, task *domain.Task) (int, error) {
	var (
		newVersion *int
		rowExists  bool
	)

	if err := r.pool.QueryRow(
		ctx,
		`WITH updated AS(
				UPDATE todo_app.tasks
					SET title = $1,
						description = $2,
						completed = $3,
						completed_at = $4,
						version = version + 1
				WHERE id = $5 AND version = $6
				RETURNING version)
		SELECT(SELECT version FROM updated),
				EXISTS(SELECT 1 FROM todo_app.tasks WHERE id = $5)`,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.ID,
		task.Version).Scan(&newVersion, &rowExists); err != nil {
		return 0, fmt.Errorf("query row: %w", utils.TranslateError(err))
	}

	if err := utils.CheckVersionedUpdate(newVersion, rowExists, domain.ErrTaskNotFound); err != nil {
		return 0, fmt.Errorf("check versioned: %w", err)
	}

	return *newVersion, nil
}
