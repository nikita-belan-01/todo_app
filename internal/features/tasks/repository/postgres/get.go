package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r TasksRepository) SearchTasks(ctx context.Context, tasksFilter *domain.TasksFilter) ([]domain.Task, error) {
	offset := utils.CalculateOffset(tasksFilter.Pagination.Page, tasksFilter.Pagination.Limit)

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`SELECT
									id,
									user_id,
									title,
									description,
									created_at,
									completed,
									completed_at,
									version
							FROM todo_app.tasks`)

	args := []any{tasksFilter.Pagination.Limit, offset}
	var conditions []string

	if tasksFilter.UserIDs != nil {
		args = append(args, tasksFilter.UserIDs)
		conditions = append(conditions, fmt.Sprintf("user_id = ANY($%d)", len(args)))
	}

	if tasksFilter.IsCompleted != nil {
		args = append(args, tasksFilter.IsCompleted)
		conditions = append(conditions, fmt.Sprintf("completed = $%d", len(args)))
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, "\nAND "))
	}

	queryBuilder.WriteString("\nORDER BY created_at, id LIMIT $1 OFFSET $2")

	rows, err := r.pool.Query(
		ctx,
		queryBuilder.String(),
		args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", utils.TranslateError(err))
	}
	defer rows.Close()

	var tasks []domain.Task

	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.CreatedAt,
			&task.Completed,
			&task.CompletedAt,
			&task.Version); err != nil {
			return nil, fmt.Errorf("scan row: %w", utils.TranslateError(err))
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", utils.TranslateError(err))
	}

	return tasks, nil
}

func (r TasksRepository) GetTask(ctx context.Context, taskID uuid.UUID) (*domain.Task, error) {
	var task domain.Task

	if err := r.pool.QueryRow(
		ctx,
		`SELECT
				user_id,
				title,
				description,
				created_at,
				completed,
				completed_at,
				version
		FROM todo_app.tasks
		WHERE id = $1`, taskID).Scan(
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.Completed,
		&task.CompletedAt,
		&task.Version); err != nil {
		if utils.IsNoRows(err) {
			return nil, domain.NewNotFoundError(domain.ErrTaskNotFound)
		}

		return nil, fmt.Errorf("scan row: %w", utils.TranslateError(err))
	}

	task.ID = taskID

	return &task, nil
}
