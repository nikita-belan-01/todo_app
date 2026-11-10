package domain

import (
	"time"

	"github.com/google/uuid"
)

type StatisticsFilter struct {
	UserID uuid.UUID
	From   *time.Time
	To     *time.Time
}

func (f StatisticsFilter) Validate() error {
	var v Validation

	if f.From != nil && f.To != nil && f.To.Compare(*f.From) <= 0 {
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
