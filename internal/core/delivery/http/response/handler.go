package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
	"go.uber.org/zap"
)

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
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)
	if err == nil {
		return
	}

	var cErr *domain.ClientError
	var msg string
	if errors.As(err, &cErr) {
		statusCode = cErr.StatusCode()
		msg = cErr.Error()
	} else {
		statusCode = http.StatusInternalServerError
		msg = "internal server error"
	}

	switch statusCode {
	case http.StatusBadRequest:
		logFunc = h.log.Warn
	case http.StatusNotFound:
		logFunc = h.log.Debug
	case http.StatusConflict:
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
