package domain

import "errors"

var (
	ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrTaskNotFound             = errors.New("task not found")
	ErrInvalidArgument          = errors.New("invalid argument")
	ErrVersionConflict          = errors.New("version conflict")
	ErrUserIDIsNull             = errors.New("user id is null")
	ErrTooManyUserIDs           = errors.New("too many user ids")
	ErrEmptyUserIDs             = errors.New("user ids is empty")
	ErrNullNotAllowed           = errors.New("null is not allowed")
	ErrInvalidUUID              = errors.New("invalid uuid")
	ErrVersionRequired          = errors.New("version is required")
	ErrEmptyPatch               = errors.New("patch is empty")
	ErrCompletedAtInconsistent  = errors.New("completed and completed_at are inconsistent")
	ErrInvalidPagination        = errors.New("invalid pagination")
	ErrInvalidDateRange         = errors.New("invalid date range")
)
