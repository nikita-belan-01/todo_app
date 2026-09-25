package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r UsersRepository) GetUsers(ctx context.Context, page, limit int) ([]domain.User, error) {
	offset := utils.CalculateOffset(page, limit)

	rows, err := r.pool.Query(
		ctx,
		`SELECT id, name, surname, phone_number FROM todo_app.users
		LIMIT $1
		OFFSET $2`,
		limit,
		offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.PhoneNumber); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r UsersRepository) GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User

	if err := r.pool.QueryRow(
		ctx,
		`SELECT id, name, surname, phone_number, version
			FROM todo_app.users
			WHERE id=$1`, userID).
		Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.PhoneNumber,
			&user.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("user not found", err)
		}

		return nil, fmt.Errorf("scan row: %w", err)
	}

	return &user, nil
}
