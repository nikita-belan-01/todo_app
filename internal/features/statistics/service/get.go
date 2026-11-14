package service

import (
	"context"
	"fmt"
	"time"

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

	return calculateStatistics(tasks), nil
}

func calculateStatistics(tasks []domain.Task) *domain.Statistics {
	if len(tasks) == 0 {
		return new(domain.Statistics)
	}

	tasksCreated := len(tasks)
	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)

		tasksAverageCompletionTime = &avg
	}

	return &domain.Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         &tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
