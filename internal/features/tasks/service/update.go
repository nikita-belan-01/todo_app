package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s TasksService) PatchTask(ctx context.Context, taskID uuid.UUID, patch *domain.TaskNullable) (int, error) {
	if err := patch.Validate(); err != nil {
		return 0, fmt.Errorf("validate: %w", err)
	}

	origTask, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("get by id: %w", err)
	}

	if origTask.Version != patch.Version {
		return 0, domain.NewAppError(domain.KindConflict, "task was modified, re-read it", domain.ErrVersionConflict)
	}

	if err = origTask.CheckTransition(patch); err != nil {
		return 0, fmt.Errorf("check transition: %w", err)
	}

	patchedTask, changed := origTask.WithPatch(patch)
	if !changed {
		return origTask.Version, nil
	}

	newVersion, err := s.tasksRepository.UpdateTask(ctx, patchedTask)
	if err != nil {
		return 0, fmt.Errorf("update: %w", err)
	}

	return newVersion, nil
}
