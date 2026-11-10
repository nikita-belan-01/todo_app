package utils

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

var ErrInvalidQueryParam = errors.New("invalid query parameter")

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
			ErrInvalidQueryParam)
	}

	return value, nil
}

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	value, err := time.Parse(time.DateOnly, param)
	if err != nil {
		return nil, domain.NewAppError(
			domain.KindInvalidArgument,
			"invalid format, expected a date",
			ErrInvalidQueryParam)
	}

	return &value, nil
}

func GetUUIDQueryParam(r *http.Request, key string) (uuid.UUID, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return uuid.Nil, nil
	}

	value, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, domain.NewAppError(
			domain.KindInvalidArgument,
			"invalid format, expected a uuid",
			ErrInvalidQueryParam)
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
