package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s TasksService) SearchTasks(ctx context.Context, tasksFilter *domain.TasksFilter) ([]domain.Task, error) {
	if err := tasksFilter.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	data, err := s.tasksRepository.SearchTasks(ctx, tasksFilter)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s TasksService) GetTask(ctx context.Context, taskID uuid.UUID) (*domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return task, nil
}
