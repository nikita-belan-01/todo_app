package postgres

import "github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/pool"

type StatisticsRepository struct {
	pool pool.Pool
}

func NewStatisticsRepository(pool pool.Pool) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
