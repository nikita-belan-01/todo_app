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

func (h TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke PatchTask handler")

	taskID, err := utils.GetUUIDPathValue(r, taskIDPathKey)
	if err != nil {
		responseHandler.ErrorResponse(domain.NewAppError(domain.KindInvalidArgument, "invalid task id", err))
		return
	}

	patch, err := request.Decode[patchTaskRequest](r)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("decode patch task request: %w", err))
		return
	}

	version, err := h.tasksService.PatchTask(
		ctx, taskID, &patch)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("patch task: %w", err))
		return
	}

	responseHandler.JSONResponse(taskPatchResponse{ID: taskID, Version: version}, http.StatusOK)
}
