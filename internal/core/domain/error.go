package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Kind int

const (
	KindInternal Kind = iota
	KindUnavailable
	KindNotFound
	KindInvalidArgument
	KindUnprocessable
	KindConflict
	KindTooLarge
)

type AppError struct {
	kind     Kind
	message  string
	problems []Problem
	err      error
}

func (ce *AppError) Error() string {
	if ce.err == nil {
		return ce.message
	}

	cause := ce.err.Error()
	if cause == ce.message {
		return ce.message
	}

	if rest, ok := strings.CutPrefix(cause, ce.message+": "); ok {
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

func (ce *AppError) Problems() []Problem {
	return ce.problems
}

func NewAppError(kind Kind, message string, err error) *AppError {
	return &AppError{
		kind:    kind,
		message: message,
		err:     err}
}

func newValidationError(message string, problems []Problem, err error) *AppError {
	return &AppError{
		kind:     KindInvalidArgument,
		message:  message,
		problems: problems,
		err:      err}
}

func NewFieldError(kind Kind, field, message string, err error) *AppError {
	problem := Problem{Field: field, Message: message}

	return &AppError{
		kind:     kind,
		message:  problem.String(),
		problems: []Problem{problem},
		err:      err}
}

func KindOf(err error) Kind {
	var e *AppError
	if errors.As(err, &e) {
		return e.kind
	}

	return KindInternal
}

func NewNotFoundError(cause error) *AppError {
	return NewAppError(KindNotFound, cause.Error(), cause)
}
