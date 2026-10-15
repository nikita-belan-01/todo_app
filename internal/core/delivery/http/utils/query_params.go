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
	page := GetIntQueryParam(r, pageQueryKey, domain.PageDefaultValue)
	if page <= 0 || page > domain.PageMaxValue {
		page = domain.PageDefaultValue
	}

	limit := GetIntQueryParam(r, limitQueryKey, domain.LimitDefaultValue)
	if limit <= 0 || limit > domain.LimitMaxValue {
		limit = domain.LimitDefaultValue
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}
