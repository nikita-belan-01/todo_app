package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

func (r UsersRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if _, err := r.pool.Exec(
		ctx,
		`INSERT INTO todo_app.users (name, surname, phone_number) 
		VALUES ($1, $2, $3);`,
		user.Name, user.Surname, user.PhoneNumber); err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return domain.NewConflictError(
					"phone number already exists",
					domain.ErrPhoneNumberAlreadyExists,
				)
			case pgerrcode.CheckViolation, pgerrcode.StringDataRightTruncationDataException:
				return domain.NewBadRequestError("invalid user data", err)
			}
		}

		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}
