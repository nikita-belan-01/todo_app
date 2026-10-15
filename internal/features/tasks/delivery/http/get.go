package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/request"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/response"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/utils"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
)

func (h TasksHTTPHandler) SearchTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke SearchTasks handler")

	tasksFilter, err := request.Decode[searchTasksRequest](r)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("decode search tasks request: %w", err))
		return
	}

	tasks, err := h.tasksService.SearchTasks(ctx, &tasksFilter)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("search tasks: %w", err))
		return
	}

	responseHandler.JSONResponse(buildTasksResponse(tasks), http.StatusOK)
}

func (h TasksHTTPHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetTask handler")

	taskID, err := utils.GetUUIDPathValue(r, taskIDPathKey)
	if err != nil {
		responseHandler.ErrorResponse(domain.NewAppError(domain.KindInvalidArgument, "invalid task id", err))
		return
	}

	task, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("get task: %w", err))
		return
	}

	responseHandler.JSONResponse(buildTaskResponse(task), http.StatusOK)
}
