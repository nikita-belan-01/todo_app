package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/response"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/utils"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
)

func (h TasksHTTPHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke DeleteTask handler")

	taskID, err := utils.GetUUIDPathValue(r, taskIDPathKey)
	if err != nil {
		responseHandler.ErrorResponse(domain.NewAppError(domain.KindInvalidArgument, "invalid task id", err))
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("delete task: %w", err))
		return
	}

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}
