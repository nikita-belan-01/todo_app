package postgres

import "github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/pool"

type UsersRepository struct {
	pool pool.Pool
}

func NewUsersRepository(pool pool.Pool) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
