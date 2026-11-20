package service

import (
	"context"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasksCompletionParams(ctx context.Context, statisticsFilter *domain.StatisticsFilter) ([]domain.Task, error)
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
