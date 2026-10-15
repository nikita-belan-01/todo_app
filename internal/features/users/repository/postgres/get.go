package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/utils"
)

func (r UsersRepository) GetUsers(ctx context.Context, page, limit int) ([]domain.User, error) {
	offset := utils.CalculateOffset(page, limit)

	rows, err := r.pool.Query(
		ctx,
		`SELECT 
				id, 
				name, 
				surname, 
				phone_number,
				version 
		FROM todo_app.users
		ORDER BY created_at, id
		LIMIT $1
		OFFSET $2`,
		limit,
		offset)
	if err != nil {
		return nil, fmt.Errorf("query: %w", utils.TranslateError(err))
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.PhoneNumber,
			&user.Version); err != nil {
			return nil, fmt.Errorf("scan row: %w", utils.TranslateError(err))
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", utils.TranslateError(err))
	}

	return users, nil
}

func (r UsersRepository) GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User

	if err := r.pool.QueryRow(
		ctx,
		`SELECT 
				name, 
				surname, 
				phone_number, 
				version
		FROM todo_app.users
		WHERE id = $1`, userID).
		Scan(
			&user.Name,
			&user.Surname,
			&user.PhoneNumber,
			&user.Version); err != nil {
		return nil, fmt.Errorf("scan row: %w", utils.TranslateError(err))
	}

	user.ID = userID

	return &user, nil
}
