package postgres

import (
	"context"
	"fmt"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r UsersRepository) UpdateUser(ctx context.Context, user *domain.User) (int, error) {
	var (
		newVersion *int
		rowExists  bool
	)

	if err := r.pool.QueryRow(
		ctx,
		`WITH updated AS(
				UPDATE todo_app.users
					SET name = $1,
						surname = $2,
						phone_number = $3,
						version = version + 1
				WHERE id = $4 AND version = $5
				RETURNING version)
		SELECT(SELECT version FROM updated),
				EXISTS(SELECT 1 FROM todo_app.users WHERE id = $4)`,
		user.Name,
		user.Surname,
		user.PhoneNumber,
		user.ID,
		user.Version).Scan(&newVersion, &rowExists); err != nil {
		return 0, fmt.Errorf("query row: %w", utils.TranslateError(err))
	}

	if err := utils.CheckVersionedUpdate(newVersion, rowExists, domain.ErrUserNotFound); err != nil {
		return 0, fmt.Errorf("check versioned: %w", err)
	}

	return *newVersion, nil
}
