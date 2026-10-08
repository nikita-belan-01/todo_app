package config

import (
	"fmt"
	"time"

	"github.com/nikita-belan-01/todo_app/pkg/default_env"
)

type Handler struct {
	RequestTimeout time.Duration
}

func newHandler() (Handler, error) {
	requestTimeout, err := default_env.GetDuration("HANDLER_REQUEST_TIMEOUT", 5*time.Second)
	if err != nil {
		return Handler{}, err
	}

	if err = validateDuration(requestTimeout); err != nil {
		return Handler{}, fmt.Errorf("HANDLER_REQUEST_TIMEOUT: %w", err)
	}

	return Handler{
		RequestTimeout: requestTimeout,
	}, nil
}
