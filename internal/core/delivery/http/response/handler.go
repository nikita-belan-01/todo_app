package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
	"go.uber.org/zap"
)

const statusClientClosedRequest = 499

type HTTPResponseHandler struct {
	log *logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(log *logger.Logger, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	if responseBody == nil || statusCode == http.StatusNoContent {
		h.rw.WriteHeader(statusCode)
		return
	}

	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}
func (h HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(http.StatusInternalServerError, msg)
}

func (h HTTPResponseHandler) ErrorResponse(err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, context.Canceled):
		h.log.Debug("request canceled by client", zap.Error(err))
		h.errorResponse(statusClientClosedRequest, "client closed request")
		return
	case errors.Is(err, context.DeadlineExceeded):
		err = domain.NewAppError(domain.KindUnavailable, "request timed out", err)
	}

	statusCode, ok := statusByKind[domain.KindOf(err)]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	msg := "internal server error"

	var appErr *domain.AppError
	if errors.As(err, &appErr) && statusCode != http.StatusInternalServerError {
		msg = appErr.Message()
	}

	var logFunc func(string, ...zap.Field)
	switch {
	case statusCode == http.StatusNotFound:
		logFunc = h.log.Debug
	case statusCode >= http.StatusBadRequest &&
		statusCode < http.StatusInternalServerError:
		logFunc = h.log.Warn
	default:
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, msg)
}

func (h HTTPResponseHandler) errorResponse(statusCode int, msg string) {
	h.JSONResponse(map[string]string{"message": msg}, statusCode)
}
