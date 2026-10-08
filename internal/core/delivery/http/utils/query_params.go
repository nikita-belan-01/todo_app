package utils

import (
	"net/http"
	"strconv"
)

const (
	PageQueryKey     = "page"
	PageDefaultValue = 1
	PageMaxValue     = 1_000_000

	LimitQueryKey     = "limit"
	LimitDefaultValue = 10
	LimitMaxValue     = 100
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
	page := GetIntQueryParam(r, PageQueryKey, PageDefaultValue)
	if page <= 0 || page > PageMaxValue {
		page = PageDefaultValue
	}

	limit := GetIntQueryParam(r, LimitQueryKey, LimitDefaultValue)
	if limit <= 0 || limit > LimitMaxValue {
		limit = LimitDefaultValue
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}
