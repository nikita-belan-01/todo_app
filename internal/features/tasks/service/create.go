package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s TasksService) CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {
	if err := task.Validate(); err != nil {
		return uuid.Nil, fmt.Errorf("validate: %w", err)
	}

	id, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
