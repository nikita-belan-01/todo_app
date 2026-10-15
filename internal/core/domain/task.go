package domain

import (
	"errors"
	"fmt"
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
	var errs []error

	if tf.UserIDs != nil {
		switch {
		case len(tf.UserIDs) == 0:
			errs = append(errs, ErrEmptyUserIDs)
		case len(tf.UserIDs) > UserIDsMaxCount:
			errs = append(errs, fmt.Errorf("%w: got %d max %d", ErrTooManyUserIDs, len(tf.UserIDs), UserIDsMaxCount))
		}

		if slices.Contains(tf.UserIDs, uuid.Nil) {
			errs = append(errs, ErrUserIDIsNull)
		}
	}

	err := errors.Join(errs...)
	if err != nil {
		return err
	}

	if tf.Page <= 0 || tf.Page > PageMaxValue {
		tf.Page = PageDefaultValue
	}

	if tf.Limit <= 0 || tf.Limit > LimitMaxValue {
		tf.Limit = LimitDefaultValue
	}

	return nil
}

func (t Task) Validate() error {
	var errs []error
	if t.UserID == uuid.Nil {
		errs = append(errs, ErrUserIDIsNull)
	}

	if err := validator.ValidateLen(t.Title, 3, 100); err != nil {
		errs = append(errs, fmt.Errorf("title: %w", err))
	}

	if t.Description != nil {
		if err := validator.ValidateLen(*t.Description, 3, 1000); err != nil {
			errs = append(errs, fmt.Errorf("description(not null): %w", err))
		}
	}

	return errors.Join(errs...)
}

type TaskNullable struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
	CompletedAt Nullable[time.Time]
}

func (t TaskNullable) Validate() error {
	if !t.Title.Set && !t.Description.Set && !t.Completed.Set && !t.CompletedAt.Set {
		return fmt.Errorf("patch must contain at least one field: %w", ErrInvalidArgument)
	}

	var errs []error

	if t.Title.Set {
		if t.Title.Value == nil {
			errs = append(errs, fmt.Errorf("title can't be patched to NULL: %w", ErrInvalidArgument))
		} else if err := validator.ValidateLen(*t.Title.Value, 3, 100); err != nil {
			errs = append(errs, fmt.Errorf("title: %w", err))
		}
	}

	if t.Description.Set && t.Description.Value != nil {
		if err := validator.ValidateLen(*t.Description.Value, 3, 1000); err != nil {
			errs = append(errs, fmt.Errorf("description: %w", err))
		}
	}

	if err := t.validateCompletion(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (t TaskNullable) validateCompletion() error {
	if t.Completed.Set != t.CompletedAt.Set {
		return fmt.Errorf("completed and completed_at must be set together: %w", ErrInvalidArgument)
	}
	if !t.Completed.Set {
		return nil
	}

	if t.Completed.Value == nil {
		return fmt.Errorf("completed can't be patched to NULL: %w", ErrInvalidArgument)
	}

	completed := *t.Completed.Value
	hasCompletedAt := t.CompletedAt.Value != nil

	switch {
	case completed && !hasCompletedAt:
		return fmt.Errorf("completed_at is required when completed is true: %w", ErrInvalidArgument)
	case !completed && hasCompletedAt:
		return fmt.Errorf("completed_at must be null when completed is false: %w", ErrInvalidArgument)
	case completed && t.CompletedAt.Value.After(time.Now()):
		return fmt.Errorf("completed_at can't be in the future: %w", ErrInvalidArgument)
	}

	return nil
}

func (t Task) WithPatch(patch *TaskNullable) (*Task, error) {
	var changed bool

	if patch.Title.Set && patch.Title.Value != nil && *patch.Title.Value != t.Title {
		t.Title = *patch.Title.Value
		changed = true
	}

	if patch.Description.Set && !equalPtr(patch.Description.Value, t.Description) {
		t.Description = patch.Description.Value
		changed = true
	}

	if patch.Completed.Set && patch.Completed.Value != nil && *patch.Completed.Value != t.Completed {
		t.Completed = *patch.Completed.Value
		changed = true
	}

	if patch.CompletedAt.Set && !equalTimePtr(patch.CompletedAt.Value, t.CompletedAt) {
		t.CompletedAt = patch.CompletedAt.Value
		changed = true

		if t.CompletedAt != nil && !t.CompletedAt.Round(time.Second).After(t.CreatedAt.Round(time.Second)) {
			return nil, NewAppError(KindInvalidArgument, "completed_at must be later than created_at", ErrInvalidArgument)
		}
	}

	if !changed {
		return nil, NewAppError(KindInvalidArgument, "nothing to update", ErrNothingToUpdate)
	}

	return &t, nil
}
