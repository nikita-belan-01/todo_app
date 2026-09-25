package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          int
	UserID      uuid.UUID
	Version     int
	Title       string
	Description string
	CreatedAt   time.Time
	Completed   bool
	CompletedAt *time.Time
}
