package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error)
	SearchTasks(ctx context.Context, tasksFilter *domain.TasksFilter) ([]domain.Task, error)
	GetTask(ctx context.Context, taskID uuid.UUID) (*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) (int, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
}

func NewTasksService(tasksRepository TasksRepository) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
