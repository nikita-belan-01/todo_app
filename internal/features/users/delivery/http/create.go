package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/request"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/response"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
)

func (h UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.Config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke CreateUser handler")

	user, err := request.Decode[createUserRequest](r)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("decode create user request: %w", err))
		return
	}

	id, err := h.UsersService.CreateUser(ctx, &user)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("create user: %w", err))
		return
	}

	responseHandler.JSONResponse(userIDResponse{ID: id}, http.StatusCreated)
}
