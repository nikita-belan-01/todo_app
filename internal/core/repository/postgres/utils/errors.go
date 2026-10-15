package utils

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

type mapped struct {
	kind     domain.Kind
	message  string
	sentinel error
}

var byConstraint = map[string]mapped{
	"tasks_user_id_fkey":          {domain.KindNotFound, "user not found", domain.ErrUserNotFound},
	"users_phone_number_key":      {domain.KindConflict, "phone number already exists", domain.ErrPhoneNumberAlreadyExists},
	"tasks_title_length":          {domain.KindInvalidArgument, "title must be 3..100 characters", domain.ErrInvalidArgument},
	"tasks_description_length":    {domain.KindInvalidArgument, "description must be 3..1000 characters", domain.ErrInvalidArgument},
	"tasks_completed_consistency": {domain.KindInvalidArgument, "completed and completed_at are inconsistent", domain.ErrInvalidArgument},
	"users_name_length":           {domain.KindInvalidArgument, "name must be 3..100 characters", domain.ErrInvalidArgument},
	"users_surname_length":        {domain.KindInvalidArgument, "surname must be 3..100 characters", domain.ErrInvalidArgument},
	"users_phone_number_valid":    {domain.KindInvalidArgument, "phone number is malformed", domain.ErrInvalidArgument},
}

var byCode = map[string]mapped{
	pgerrcode.UniqueViolation:                        {domain.KindConflict, "already exists", nil},
	pgerrcode.ForeignKeyViolation:                    {domain.KindInvalidArgument, "referenced entity not found", nil},
	pgerrcode.CheckViolation:                         {domain.KindInvalidArgument, "invalid data", nil},
	pgerrcode.NotNullViolation:                       {domain.KindInvalidArgument, "required field is missing", nil},
	pgerrcode.StringDataRightTruncationDataException: {domain.KindInvalidArgument, "value is too long", nil},
	pgerrcode.QueryCanceled:                          {domain.KindUnavailable, "request timed out", nil},
}

func TranslateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewAppError(domain.KindNotFound, "not found", err)
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	m, ok := byConstraint[pgErr.ConstraintName]
	if !ok {
		if m, ok = byCode[pgErr.Code]; !ok {
			return err
		}
	}

	cause := error(pgErr)
	if m.sentinel != nil {
		cause = fmt.Errorf("%w: %w", m.sentinel, pgErr)
	}

	return domain.NewAppError(m.kind, m.message, cause)
}

func CheckVersionedUpdate(newVersion *int, rowExists bool, notFoundErr error) error {
	switch {
	case newVersion != nil:
		return nil
	case rowExists:
		return domain.NewAppError(domain.KindConflict, "version conflict", domain.ErrVersionConflict)
	default:
		return domain.NewAppError(domain.KindNotFound, "not found", notFoundErr)
	}
}
