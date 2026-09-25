package utils

import (
	"net/http"
	"strconv"
)

const (
	PageQueryKey  = "page"
	LimitQueryKey = "limit"
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

	if value == 0 {
		return defaultValue
	}

	return value
}
