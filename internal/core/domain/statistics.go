package domain

import (
	"time"

	"github.com/google/uuid"
)

type StatisticsFilter struct {
	UserID *uuid.UUID
	From   *time.Time
	To     *time.Time
}

func (f StatisticsFilter) Validate() error {
	var v Validation

	if f.UserID != nil && *f.UserID == uuid.Nil {
		v.AddWithCause("user_id", ErrUserIDIsNull, "must not be empty")
	}

	if f.From != nil && f.To != nil && f.To.Before(*f.From) {
		v.AddWithCause("to", ErrInvalidDateRange, "must not be earlier than from")
	}

	return v.Err()
}

type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func CreateStatistics(tasks []Task) *Statistics {
	if len(tasks) == 0 {
		return new(Statistics)
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

	return &Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         &tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
