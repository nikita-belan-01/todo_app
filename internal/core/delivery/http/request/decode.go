package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
)

const unknownFieldPrefix = "json: unknown field"

func Decode[T interface{ BuildDomain() DT }, DT any](r *http.Request) (DT, error) {
	var dto T
	var zero DT

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dto); err != nil {
		return zero, decodeError(err)
	}

	var extra json.RawMessage
	switch err := dec.Decode(&extra); {
	case err == nil:
		return zero, domain.NewAppError(domain.KindInvalidArgument, "unexpected content after json body", nil)
	case !errors.Is(err, io.EOF):
		return zero, decodeError(err)
	}

	return dto.BuildDomain(), nil
}

func decodeError(err error) error {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return domain.NewAppError(domain.KindTooLarge, fmt.Sprintf("request body must not exceed %d bytes", maxBytesErr.Limit), err)
	}

	if errors.Is(err, io.EOF) {
		return domain.NewAppError(domain.KindInvalidArgument, "empty request body", err)
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		return domain.NewAppError(domain.KindInvalidArgument, "malformed json body: unexpected end of input", err)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return domain.NewAppError(
			domain.KindInvalidArgument,
			fmt.Sprintf("field %q: expected %s, got %s", typeErr.Field, typeErr.Type, typeErr.Value), err)
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return domain.NewAppError(
			domain.KindInvalidArgument,
			fmt.Sprintf("malformed json body at offset %d", syntaxErr.Offset), err)
	}

	if field, ok := strings.CutPrefix(err.Error(), unknownFieldPrefix); ok {
		return domain.NewAppError(
			domain.KindInvalidArgument,
			fmt.Sprintf("unknown field %s", strings.TrimSpace(field)), err)
	}

	return domain.NewAppError(domain.KindInvalidArgument, "invalid json body request", err)
}
