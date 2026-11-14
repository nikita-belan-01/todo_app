package http

import (
	"net/http"

	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/utils"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func getStatisticsFilter(r *http.Request) (*domain.StatisticsFilter, error) {
	const (
		userIDQueryKey = "user_id"
		fromQueryKey   = "from"
		toQueryKey     = "to"
	)

	var v domain.Validation

	userID, err := utils.GetUUIDQueryParam(r, userIDQueryKey)
	if err != nil {
		v.Add(userIDQueryKey, err)
	}

	from, err := utils.GetDateQueryParam(r, fromQueryKey)
	if err != nil {
		v.Add(fromQueryKey, err)
	}

	to, err := utils.GetDateQueryParam(r, toQueryKey)
	if err != nil {
		v.Add(toQueryKey, err)
	}

	if err := v.Err(); err != nil {
		return nil, err
	}

	return &domain.StatisticsFilter{
		UserID: userID,
		From:   from,
		To:     to,
	}, nil
}
