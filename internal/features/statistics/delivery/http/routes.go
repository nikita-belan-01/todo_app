package http

import (
	"context"
	"net/http"

	"github.com/nikita-belan-01/todo_app/config"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/server"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

type StatisticsHTTPHandler struct {
	config            *config.Handler
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(ctx context.Context, statisticsFilter *domain.StatisticsFilter) (*domain.Statistics, error)
}

func NewTasksHTTPHandler(config *config.Handler, statisticsService StatisticsService) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		config:            config,
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []server.Route {
	const statisticsPath string = "/statistics"
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    statisticsPath,
			Handler: h.GetStatistics,
		}}
}
