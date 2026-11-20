package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

var (
	ErrEmptyPathValue   = errors.New("empty path value")
	ErrInvalidPathValue = errors.New("invalid path value")
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, fmt.Errorf("no key=%q in path values: %w", key, ErrEmptyPathValue)
	}

	id, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, domain.NewFieldError(
			domain.KindInvalidArgument,
			key,
			"invalid format, expected a uuid",
			fmt.Errorf("%w: %w", ErrInvalidPathValue, err))
	}

	return id, err
}
