package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/config"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/server"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

const taskIDPathKey = "task_id"

type TasksHTTPHandler struct {
	config       *config.Handler
	tasksService TasksService
}

type TasksService interface {
	CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error)
	SearchTasks(ctx context.Context, tasksFilter *domain.TasksFilter) ([]domain.Task, error)
	GetTask(ctx context.Context, taskID uuid.UUID) (*domain.Task, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	PatchTask(ctx context.Context, taskID uuid.UUID, patch *domain.TaskNullable) (int, error)
}

func NewTasksHTTPHandler(config *config.Handler, tasksService TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		config:       config,
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []server.Route {
	const tasksPath string = "/tasks"
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    tasksPath,
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodPost,
			Path:    fmt.Sprintf("%s/search", tasksPath),
			Handler: h.SearchTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    fmt.Sprintf("%s/{%s}", tasksPath, taskIDPathKey),
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    fmt.Sprintf("%s/{%s}", tasksPath, taskIDPathKey),
			Handler: h.DeleteTask,
		}, {
			Method:  http.MethodPatch,
			Path:    fmt.Sprintf("%s/{%s}", tasksPath, taskIDPathKey),
			Handler: h.PatchTask,
		}}
}
