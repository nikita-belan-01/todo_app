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

func (r UsersRepository) UpdateUser(ctx context.Context, userID uuid.UUID, user *domain.User) error {
	cmd, err := r.pool.Exec(
		ctx,
		`UPDATE todo_app.users
			SET name = $1,
				surname = $2,
				phone_number = $3,
				version = version + 1
			WHERE id = $4 AND version = $5;`,
		user.Name, user.Surname, user.PhoneNumber, userID, user.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return domain.NewConflictError(
					"phone number already exists", domain.ErrPhoneNumberAlreadyExists)
			case pgerrcode.CheckViolation, pgerrcode.StringDataRightTruncationDataException:
				return domain.NewBadRequestError("invalid user data", err)
			}
		}

		return fmt.Errorf("update user: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return domain.NewConflictError("conflict user version", domain.ErrVersionConflict)
	}

	return nil
}
