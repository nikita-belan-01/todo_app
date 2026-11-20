package domain

import (
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/pkg/validator"
)

type Task struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Version     int
	Title       string
	Description *string
	CreatedAt   time.Time
	Completed   bool
	CompletedAt *time.Time
}

func (t Task) CompletionDuration() *time.Duration {
	if !t.Completed || t.CompletedAt == nil {
		return nil
	}

	duration := t.CompletedAt.Sub(t.CreatedAt)

	return &duration
}

type TasksFilter struct {
	UserIDs     uuid.UUIDs
	IsCompleted *bool
	Pagination  Pagination
}

func (tf *TasksFilter) Validate() error {
	var v Validation

	if tf.UserIDs != nil {
		switch {
		case len(tf.UserIDs) == 0:
			v.AddWithCause("user_ids", ErrEmptyUserIDs, "must not be empty")
		case len(tf.UserIDs) > UserIDsMaxCount:
			v.AddWithCause("user_ids", ErrTooManyUserIDs, "must contain at most %d ids, got %d", UserIDsMaxCount, len(tf.UserIDs))
		}

		if slices.Contains(tf.UserIDs, uuid.Nil) {
			v.AddWithCause("user_ids", ErrUserIDIsNull, "must not contain an empty id")
		}
	}

	for _, err := range tf.Pagination.validate() {
		v.Add(err.field, err)
	}

	return v.Err()
}

func (t Task) Validate() error {
	var v Validation

	if t.UserID == uuid.Nil {
		v.AddWithCause("user_id", ErrUserIDIsNull, "must not be empty")
	}

	v.Add("title", validator.ValidateLen(t.Title, TitleMinLen, TitleMaxLen))

	if t.Description != nil {
		v.Add("description", validator.ValidateLen(*t.Description, DescriptionMinLen, DescriptionMaxLen))
	}

	switch {
	case !t.Completed && t.CompletedAt != nil:
		v.AddWithCause("completed_at", ErrCompletedAtInconsistent, "must be null when completed is false")
	case t.Completed && t.CompletedAt == nil:
		v.AddWithCause("completed_at", ErrCompletedAtInconsistent, "is required when completed is true")
	case t.Completed && t.CompletedAt.Before(t.CreatedAt):
		v.AddWithCause("completed_at", ErrCompletedAtInconsistent, "must not be earlier than created_at")
	}

	return v.Err()
}

type TaskNullable struct {
	Version     int
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func (t TaskNullable) Validate() error {
	var v Validation

	if t.Version < 1 {
		v.AddWithCause("version", ErrVersionRequired, "is required and must be greater than 0")
	}

	if !t.Title.Set && !t.Description.Set && !t.Completed.Set {
		v.AddWithCause("", ErrEmptyPatch, "patch must contain at least one field")
		return v.Err()
	}

	v.AddNullableLen("title", t.Title, TitleMinLen, TitleMaxLen)
	v.AddOptionalLen("description", t.Description, DescriptionMinLen, DescriptionMaxLen)
	v.AddNullableBool("completed", t.Completed)

	return v.Err()
}

func (t Task) WithPatch(patch *TaskNullable) (*Task, bool) {
	t.Version = patch.Version
	changed := applyNotNull(&t.Title, patch.Title)
	changed = applyNullable(&t.Description, patch.Description) || changed
	changed = applyNotNull(&t.Completed, patch.Completed) || changed

	return &t, changed
}
