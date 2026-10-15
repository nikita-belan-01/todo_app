package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s TasksService) CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {
	if err := task.Validate(); err != nil {
		return uuid.Nil, domain.NewAppError(domain.KindInvalidArgument, "invalid data", err)
	}

	id, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
