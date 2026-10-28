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

type TasksFilter struct {
	UserIDs     uuid.UUIDs
	IsCompleted *bool
	Page        int
	Limit       int
}

func (tf *TasksFilter) Validate() error {
	var v Validation

	if tf.UserIDs != nil {
		switch {
		case len(tf.UserIDs) == 0:
			v.Addc("user_ids", ErrEmptyUserIDs, "must not be empty")
		case len(tf.UserIDs) > UserIDsMaxCount:
			v.Addc("user_ids", ErrTooManyUserIDs, "must contain at most %d ids, got %d", UserIDsMaxCount, len(tf.UserIDs))
		}

		if slices.Contains(tf.UserIDs, uuid.Nil) {
			v.Addc("user_ids", ErrUserIDIsNull, "must not contain an empty id")
		}
	}

	return v.Err()
}

func (t Task) Validate() error {
	var v Validation

	if t.UserID == uuid.Nil {
		v.Addc("user_id", ErrUserIDIsNull, "must not be empty")
	}

	v.Add("title", validator.ValidateLen(t.Title, TitleMinLen, TitleMaxLen))

	if t.Description != nil {
		v.Add("description", validator.ValidateLen(*t.Description, DescriptionMinLen, DescriptionMaxLen))
	}

	switch {
	case !t.Completed && t.CompletedAt != nil:
		v.Addf("completed_at", "must be null when completed is false")
	case t.Completed && t.CompletedAt == nil:
		v.Addf("completed_at", "is required when completed is true")
	case t.Completed && t.CompletedAt.Before(t.CreatedAt):
		v.Addf("completed_at", "must not be earlier than created_at")
	}

	return v.Err()
}

type TaskNullable struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func (t TaskNullable) Validate() error {
	var v Validation

	if !t.Title.Set && !t.Description.Set && !t.Completed.Set {
		v.Addf("", "patch must contain at least one field")
		return v.Err()
	}

	v.AddNullableLen("title", t.Title, TitleMinLen, TitleMaxLen)
	v.AddOptionalLen("description", t.Description, DescriptionMinLen, DescriptionMaxLen)
	v.AddNullableBool("completed", t.Completed)

	return v.Err()
}

func (t Task) WithPatch(patch *TaskNullable) (*Task, bool) {
	changed := applyNotNull(&t.Title, patch.Title)
	changed = applyNullable(&t.Description, patch.Description) || changed
	changed = applyNotNull(&t.Completed, patch.Completed) || changed

	return &t, changed
}

func (t Task) CheckTransition(patch *TaskNullable) error {
	if !t.Completed || !patch.Completed.Set || patch.Completed.Value == nil {
		return nil
	}

	if !*patch.Completed.Value {
		return NewAppError(KindConflict, "a completed task cannot be reopened", ErrTaskCompleted)
	}

	return nil
}
