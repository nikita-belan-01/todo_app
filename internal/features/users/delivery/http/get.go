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

func (h UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetUsers handler")

	pagination, err := utils.GetPagination(r)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("pagination: %w", err))
		return
	}

	data, err := h.usersService.GetUsers(ctx, &pagination)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("get users: %w", err))
		return
	}

	responseHandler.JSONResponse(BuildUsersResponse(data), http.StatusOK)
}

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetUser handler")

	userID, err := utils.GetUUIDPathValue(r, userIDPathKey)
	if err != nil {
		responseHandler.ErrorResponse(domain.NewAppError(domain.KindInvalidArgument, "invalid user id", err))
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("get user: %w", err))
		return
	}

	responseHandler.JSONResponse(BuildUserResponse(user), http.StatusOK)
}
