package domain

import (
	"errors"
	"net/http"
)

var (
	ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrInvalidArgument          = errors.New("invalid argument")
	ErrVersionConflict          = errors.New("version conflict")
)

type ClientError struct {
	statusCode int
	message    string
	err        error
}

func NewBadRequestError(message string, err error) *ClientError {
	return &ClientError{
		statusCode: http.StatusBadRequest,
		message:    message,
		err:        err,
	}
}

func NewConflictError(message string, err error) *ClientError {
	return &ClientError{
		statusCode: http.StatusConflict,
		message:    message,
		err:        err,
	}
}

func NewNotFoundError(message string, err error) *ClientError {
	return &ClientError{
		statusCode: http.StatusNotFound,
		message:    message,
		err:        err,
	}
}

func (ce *ClientError) Error() string {
	return ce.message
}

func (ce *ClientError) Unwrap() error {
	return ce.err
}

func (ce *ClientError) StatusCode() int {
	return ce.statusCode
}
