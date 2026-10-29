package http

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/request"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/types"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/pkg/strutil"
)

const taskIDPathKey = "taskId"

type createTaskRequest struct {
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

func (r createTaskRequest) BuildDomain() (domain.Task, error) {
	var v domain.Validation

	userID, err := uuid.Parse(r.UserID)
	if err != nil && r.UserID != "" {
		v.Addf("user_id", "%q is not a valid uuid", strutil.Shorten(r.UserID, request.MaxEchoedValue))
	}

	if err := v.Err(); err != nil {
		return domain.Task{}, err
	}

	return domain.Task{
		UserID:      userID,
		Title:       r.Title,
		Description: r.Description,
	}, nil
}

type taskIDResponse struct {
	ID uuid.UUID `json:"id"`
}

type searchTasksRequest struct {
	UserIDs     []string `json:"user_ids"`
	Page        *int     `json:"page"`
	Limit       *int     `json:"limit"`
	IsCompleted *bool    `json:"is_completed"`
}

func (r searchTasksRequest) BuildDomain() (domain.TasksFilter, error) {
	var (
		v       domain.Validation
		userIDs uuid.UUIDs
	)

	if r.UserIDs != nil {
		if len(r.UserIDs) > domain.UserIDsMaxCount {
			return domain.TasksFilter{}, domain.NewAppError(
				domain.KindInvalidArgument,
				fmt.Sprintf("user_ids: must contain at most %d ids, got %d", domain.UserIDsMaxCount, len(r.UserIDs)),
				domain.ErrTooManyUserIDs)
		}
		userIDs = make(uuid.UUIDs, 0, len(r.UserIDs))

		for i, value := range r.UserIDs {
			parsed, err := uuid.Parse(value)
			if err != nil && value != "" {
				v.Addf(fmt.Sprintf("user_ids.%d", i), "%q is not a valid uuid", strutil.Shorten(value, request.MaxEchoedValue))
			}

			userIDs = append(userIDs, parsed)
		}
	}

	if err := v.Err(); err != nil {
		return domain.TasksFilter{}, err
	}

	filter := domain.TasksFilter{
		UserIDs:     userIDs,
		IsCompleted: r.IsCompleted,
		Page:        domain.PageDefaultValue,
		Limit:       domain.LimitDefaultValue,
	}

	if r.Page != nil {
		filter.Page = *r.Page
	}

	if r.Limit != nil {
		filter.Limit = *r.Limit
	}

	return filter, nil
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
	Version     int                    `json:"version"`
	Title       types.Nullable[string] `json:"title"`
	Description types.Nullable[string] `json:"description"`
	Completed   types.Nullable[bool]   `json:"completed"`
}

func (r patchTaskRequest) BuildDomain() (domain.TaskNullable, error) {
	return domain.TaskNullable{
		Version:     r.Version,
		Title:       r.Title.Nullable,
		Description: r.Description.Nullable,
		Completed:   r.Completed.Nullable,
	}, nil
}

type taskPatchResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`
}
