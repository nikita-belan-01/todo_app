package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (r UsersRepository) CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	var id uuid.UUID
	if err := r.pool.QueryRow(ctx,
		`INSERT INTO todo_app.users (
									name, 
									surname,
									phone_number) 
		VALUES ($1, $2, $3) RETURNING id`,
		user.Name,
		user.Surname,
		user.PhoneNumber).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return uuid.Nil, domain.NewConflictError(
					"phone number already exists",
					domain.ErrPhoneNumberAlreadyExists,
				)
			case pgerrcode.CheckViolation, pgerrcode.StringDataRightTruncationDataException:
				return uuid.Nil, domain.NewBadRequestError("invalid user data", err)
			}
		}

		return uuid.Nil, fmt.Errorf("insert user: %w", err)
	}

	return id, nil
}
