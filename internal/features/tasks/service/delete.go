package service

import (
	"context"

	"github.com/google/uuid"
)

func (s TasksService) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	if err := s.tasksRepository.DeleteTask(ctx, taskID); err != nil {
		return err
	}

	return nil
}
