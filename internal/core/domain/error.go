package domain

import (
	"errors"
	"fmt"
	"strings"
)

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
	ErrTaskCompleted            = errors.New("task completed")
)

type Kind int

const (
	KindInternal Kind = iota
	KindUnavailable
	KindNotFound
	KindInvalidArgument
	KindConflict
	KindTooLarge
)

type AppError struct {
	kind    Kind
	message string
	err     error
}

func (ce *AppError) Error() string {
	if ce.err == nil {
		return ce.message
	}

	cause := ce.err.Error()
	if cause == ce.message {
		return ce.message
	}

	if rest, ok := strings.CutPrefix(cause, fmt.Sprintf("%s: ", ce.message)); ok {
		cause = rest
	}

	return fmt.Sprintf("%s: %s", ce.message, cause)
}

func (ce *AppError) Unwrap() error {
	return ce.err
}

func (ce *AppError) Message() string {
	return ce.message
}

func NewAppError(kind Kind, message string, err error) *AppError {
	return &AppError{kind: kind, message: message, err: err}
}

func KindOf(err error) Kind {
	var e *AppError
	if errors.As(err, &e) {
		return e.kind
	}

	return KindInternal
}
