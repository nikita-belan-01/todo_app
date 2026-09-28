package request

import (
	"encoding/json"
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func Decode[T interface{ BuildDomain() DT }, DT any](r *http.Request) (DT, error) {
	var dto T
	var zero DT

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return zero, domain.NewBadRequestError("invalid json body request", err)
	}

	return dto.BuildDomain(), nil
}
