package request

import (
	"encoding/json"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func DecodeAndValidate[T interface{ Validate() error }](r *http.Request) (T, error) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, domain.NewBadRequestError("invalid json body request", err)
	}

	if err := req.Validate(); err != nil {
		return req, domain.NewBadRequestError("validation failed", err)
	}

	return req, nil
}
