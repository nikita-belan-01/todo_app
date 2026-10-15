package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s TasksService) PatchTask(ctx context.Context, taskID uuid.UUID, patch *domain.TaskNullable) (int, error) {
	if err := patch.Validate(); err != nil {
		return 0, domain.NewAppError(domain.KindInvalidArgument, "invalid data", err)
	}

	origTask, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("get by id: %w", err)
	}

	patchedTask, err := origTask.WithPatch(patch)
	if err != nil {
		return 0, fmt.Errorf("patching: %w", err)
	}

	newVersion, err := s.tasksRepository.UpdateTask(ctx, patchedTask)
	if err != nil {
		return 0, fmt.Errorf("update: %w", err)
	}

	return newVersion, nil
}
