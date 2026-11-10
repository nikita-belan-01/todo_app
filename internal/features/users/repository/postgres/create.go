package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r UsersRepository) CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	var id uuid.UUID
	if err := r.pool.QueryRow(
		ctx,
		`INSERT INTO todo_app.users (
									name, 
									surname,
									phone_number) 
		VALUES ($1, $2, $3) RETURNING id`,
		user.Name,
		user.Surname,
		user.PhoneNumber).Scan(&id); err != nil {

		return uuid.Nil, fmt.Errorf("query row: %w", utils.TranslateError(err))
	}

	return id, nil
}
