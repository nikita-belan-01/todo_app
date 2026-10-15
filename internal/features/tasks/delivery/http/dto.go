package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/types"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

const taskIDPathKey = "taskId"

type createTaskRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
}

func (r createTaskRequest) BuildDomain() domain.Task {
	return domain.Task{
		UserID:      r.UserID,
		Title:       r.Title,
		Description: r.Description,
	}
}

type taskIDResponse struct {
	ID uuid.UUID `json:"id"`
}

type searchTasksRequest struct {
	UserIDs     uuid.UUIDs `json:"user_ids"`
	IsCompleted *bool      `json:"is_completed"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
}

func (r searchTasksRequest) BuildDomain() domain.TasksFilter {
	return domain.TasksFilter{
		UserIDs:     r.UserIDs,
		IsCompleted: r.IsCompleted,
		Page:        r.Page,
		Limit:       r.Limit,
	}
}

type taskResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
	Version     int        `json:"version"`
}

func buildTaskResponse(task *domain.Task) taskResponse {
	return taskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		Completed:   task.Completed,
		CompletedAt: task.CompletedAt,
		Version:     task.Version,
	}
}

func buildTasksResponse(tasks []domain.Task) []taskResponse {
	response := make([]taskResponse, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, buildTaskResponse(&task))
	}

	return response
}

type patchTaskRequest struct {
	Title       types.Nullable[string] `json:"title"`
	Description types.Nullable[string] `json:"description"`
	Completed   types.Nullable[bool]   `json:"completed"`
}

func (r patchTaskRequest) BuildDomain() domain.TaskNullable {
	return domain.TaskNullable{
		Title:       r.Title.Nullable,
		Description: r.Description.Nullable,
		Completed:   r.Completed.Nullable,
	}
}

type taskPatchResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`
}
