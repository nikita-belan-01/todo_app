package domain

import (
	"errors"
	"fmt"
)

var (
	ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrTaskNotFound             = errors.New("task not found")
	ErrInvalidArgument          = errors.New("invalid argument")
	ErrVersionConflict          = errors.New("version conflict")
	ErrNothingToUpdate          = errors.New("nothing to update")
	ErrUserIDIsNull             = errors.New("user id is null")
	ErrTooManyUserIDs           = errors.New("too many user ids")
	ErrEmptyUserIDs             = errors.New("user ids is empty")
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
	if ce.err != nil {
		return fmt.Sprintf("%s: %s", ce.message, ce.err.Error())
	}

	return ce.message
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
