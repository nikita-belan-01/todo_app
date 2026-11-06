package postgres

import "github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/pool"

type TasksRepository struct {
	pool pool.Pool
}

func NewTasksRepository(pool pool.Pool) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
