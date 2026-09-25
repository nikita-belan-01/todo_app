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

func (h UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.Config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke PatchUser handler")

	userID, err := utils.GetUUIDPathValue(r, userIDPathKey)
	if err != nil {
		responseHandler.ErrorResponse(domain.NewBadRequestError("invalid user id", err))
		return
	}

	req, err := request.DecodeAndValidate[patchUserRequest](r)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("decode and validate patch user request: %w", err))
		return
	}

	if err := h.UsersService.PatchUser(
		ctx, userID, BuildDomainNullableUser(req.Name, req.Surname, req.PhoneNumber)); err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("patch user: %w", err))
		return
	}

	responseHandler.JSONResponse(nil, http.StatusOK)
}
