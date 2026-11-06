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

func GetIntQueryParam(r *http.Request, key string, defaultValue int) (int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, domain.NewAppError(
			domain.KindInvalidArgument,
			"invalid format, expected an integer",
			domain.ErrInvalidArgument)
	}

	return value, nil
}

func GetPagination(r *http.Request) (domain.Pagination, error) {
	var v domain.Validation

	page, err := GetIntQueryParam(r, pageQueryKey, domain.PageDefaultValue)
	if err != nil {
		v.Add(pageQueryKey, err)
	}

	limit, err := GetIntQueryParam(r, limitQueryKey, domain.LimitDefaultValue)
	if err != nil {
		v.Add(limitQueryKey, err)
	}

	if err := v.Err(); err != nil {
		return domain.Pagination{}, err
	}

	return domain.Pagination{
		Page:  page,
		Limit: limit,
	}, nil
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
