package utils

import (
	"net/http"
	"strconv"
)

const (
	PageQueryKey      = "page"
	PageDefaultValue  = 1
	LimitQueryKey     = "limit"
	LimitDefaultValue = 10
)

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

func GetPageQueryParam(r *http.Request) int {
	page := GetIntQueryParam(r, PageQueryKey, PageDefaultValue)
	if page <= 0 {
		return PageDefaultValue
	}

	return page
}

func GetLimitQueryParam(r *http.Request) int {
	limit := GetIntQueryParam(r, PageQueryKey, LimitDefaultValue)
	if limit <= 0 {
		return LimitDefaultValue
	}

	return limit
}
