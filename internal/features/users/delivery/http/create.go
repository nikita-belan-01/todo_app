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

	req, err := request.DecodeAndValidate[createUserRequest](r)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("decode and validate create user request: %w", err))
		return
	}

	if err := h.UsersService.CreateUser(ctx,
		BuildDomainUser(req.Name, req.Surname, req.PhoneNumber)); err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("create user: %w", err))
		return
	}

	responseHandler.JSONResponse(nil, http.StatusCreated)
}
