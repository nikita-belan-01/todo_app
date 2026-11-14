package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r StatisticsRepository) GetTasksCompletionParams(ctx context.Context, statisticsFilter *domain.StatisticsFilter) ([]domain.Task, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT 
									completed,
									created_at,
									completed_at
							FROM todo_app.tasks`)

	var args []any
	var conditions []string

	if statisticsFilter.UserID != uuid.Nil {
		args = append(args, statisticsFilter.UserID)
		conditions = append(conditions, fmt.Sprintf("user_id=$%d", len(args)))
	}

	if statisticsFilter.From != nil {
		args = append(args, statisticsFilter.From)
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", len(args)))
	}

	if statisticsFilter.To != nil {
		args = append(args, statisticsFilter.To)
		conditions = append(conditions, fmt.Sprintf("created_at<$%d", len(args)))
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, "\nAND "))
	}

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", utils.TranslateError(err))
	}
	defer rows.Close()

	var tasks []domain.Task

	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt); err != nil {
			return nil, fmt.Errorf("scan row: %w", utils.TranslateError(err))
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", utils.TranslateError(err))
	}

	return tasks, nil
}
