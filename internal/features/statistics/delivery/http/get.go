package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/response"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
)

func (h StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetStatistics handler")

	statisticsFilter, err := getStatisticsFilter(r)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("statistics filter: %w", err))
		return
	}

	data, err := h.statisticsService.GetStatistics(ctx, statisticsFilter)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("get statistics: %w", err))
		return
	}

	responseHandler.JSONResponse(buildStatisticsResponse(data), http.StatusOK)
}
