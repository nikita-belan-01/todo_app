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
	field    string
	message  string
	sentinel error
}

var byConstraint = map[string]mapped{
	"tasks_user_id_fkey":          {domain.KindUnprocessable, "user_id", "referenced user does not exist", domain.ErrUserNotFound},
	"users_phone_number_key":      {domain.KindConflict, "phone_number", "already exists", domain.ErrPhoneNumberAlreadyExists},
	"tasks_title_length":          {domain.KindInvalidArgument, "title", "length must be between 3 and 100 characters", domain.ErrInvalidArgument},
	"tasks_description_length":    {domain.KindInvalidArgument, "description", "length must be between 3 and 1000 characters", domain.ErrInvalidArgument},
	"tasks_completed_consistency": {domain.KindInvalidArgument, "completed_at", "is inconsistent with completed", domain.ErrCompletedAtInconsistent},
	"users_name_length":           {domain.KindInvalidArgument, "name", "length must be between 3 and 100 characters", domain.ErrInvalidArgument},
	"users_surname_length":        {domain.KindInvalidArgument, "surname", "length must be between 3 and 100 characters", domain.ErrInvalidArgument},
	"users_phone_number_valid":    {domain.KindInvalidArgument, "phone_number", "invalid format, expected a leading + followed by digits only", domain.ErrInvalidArgument},
}

var byCode = map[string]mapped{
	pgerrcode.UniqueViolation:                        {domain.KindConflict, "", "already exists", nil},
	pgerrcode.ForeignKeyViolation:                    {domain.KindUnprocessable, "", "referenced entity does not exist", nil},
	pgerrcode.CheckViolation:                         {domain.KindInvalidArgument, "", "request violates a data constraint", nil},
	pgerrcode.NotNullViolation:                       {domain.KindInvalidArgument, "", "required field is missing", nil},
	pgerrcode.StringDataRightTruncationDataException: {domain.KindInvalidArgument, "", "value is too long", nil},
	pgerrcode.QueryCanceled:                          {domain.KindUnavailable, "", "request timed out", nil},
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

	if m.field == "" {
		return domain.NewAppError(m.kind, m.message, cause)
	}

	return domain.NewFieldError(m.kind, m.field, m.message, cause)
}

func CheckVersionedUpdate(newVersion *int, rowExists bool, notFoundErr error) error {
	switch {
	case newVersion != nil:
		return nil
	case rowExists:
		return domain.NewAppError(domain.KindConflict, "version conflict", domain.ErrVersionConflict)
	default:
		return domain.NewNotFoundError(notFoundErr)
	}
}

func IsNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
