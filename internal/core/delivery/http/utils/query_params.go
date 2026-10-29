package utils

import (
	"net/http"
	"strconv"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

const (
	pageQueryKey  = "page"
	limitQueryKey = "limit"
)

type Pagination struct {
	Page  int
	Limit int
}

func GetIntQueryParam(r *http.Request, key string, defaultValue int) int {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}

	return value
}

func GetPagination(r *http.Request) Pagination {
	page, limit := domain.ClampPagination(
		GetIntQueryParam(r, pageQueryKey, domain.PageDefaultValue),
		GetIntQueryParam(r, limitQueryKey, domain.LimitDefaultValue))

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}

func RejectQueryParams(r *http.Request, message string) error {
	if r.URL.RawQuery != "" {
		return domain.NewAppError(
			domain.KindInvalidArgument,
			message,
			domain.ErrInvalidArgument)
	}

	return nil
}
