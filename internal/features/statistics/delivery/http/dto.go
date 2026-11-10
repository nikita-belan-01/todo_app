package http

import (
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

type statisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

func buildStatisticsResponse(statistics *domain.Statistics) statisticsResponse {
	response := statisticsResponse{
		TasksCreated:       statistics.TasksCreated,
		TasksCompleted:     statistics.TasksCompleted,
		TasksCompletedRate: statistics.TasksCompletedRate}

	if statistics.TasksAverageCompletionTime != nil {
		averageCompletionTime := statistics.TasksAverageCompletionTime.String()
		response.TasksAverageCompletionTime = &averageCompletionTime
	}

	return response
}
