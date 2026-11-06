package response

import (
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

const (
	statusClientClosedRequest = 499
	statusCodeUninitialized   = -1
)

var statusByKind = map[domain.Kind]int{
	domain.KindInternal:        http.StatusInternalServerError,
	domain.KindUnavailable:     http.StatusServiceUnavailable,
	domain.KindNotFound:        http.StatusNotFound,
	domain.KindInvalidArgument: http.StatusBadRequest,
	domain.KindUnprocessable:   http.StatusUnprocessableEntity,
	domain.KindConflict:        http.StatusConflict,
	domain.KindTooLarge:        http.StatusRequestEntityTooLarge}

func (h HTTPResponseHandler) errorResponse(statusCode int, msg string, problems []domain.Problem) {
	h.JSONResponse(errorBody{Message: msg, Problems: buildProblems(problems)}, statusCode)
}
