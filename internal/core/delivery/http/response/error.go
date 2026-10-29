package response

import (
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

var statusByKind = map[domain.Kind]int{
	domain.KindInternal:        http.StatusInternalServerError,
	domain.KindUnavailable:     http.StatusServiceUnavailable,
	domain.KindNotFound:        http.StatusNotFound,
	domain.KindInvalidArgument: http.StatusBadRequest,
	domain.KindUnprocessable:   http.StatusUnprocessableEntity,
	domain.KindConflict:        http.StatusConflict,
	domain.KindTooLarge:        http.StatusRequestEntityTooLarge}
