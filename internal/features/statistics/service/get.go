package service

import (
	"context"
	"fmt"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, statisticsFilter *domain.StatisticsFilter) (*domain.Statistics, error) {
	if err := statisticsFilter.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	tasks, err := s.statisticsRepository.GetTasksCompletionParams(ctx, statisticsFilter)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}

	return domain.CreateStatistics(tasks), nil
}
