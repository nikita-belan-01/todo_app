package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var ErrEmptyPathValue = errors.New("empty path value")

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, fmt.Errorf("no key='%s' in path values: %w", key, ErrEmptyPathValue)
	}

	id, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf("key='%s' in path values: %w", key, err)
	}

	return id, err
}
