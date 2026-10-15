package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/request"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/response"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
)

func (h TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke CreateTask handler")

	task, err := request.Decode[createTaskRequest](r)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("decode create task request: %w", err))
		return
	}

	id, err := h.tasksService.CreateTask(ctx, &task)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("create task: %w", err))
		return
	}

	responseHandler.JSONResponse(taskIDResponse{ID: id}, http.StatusCreated)
}
