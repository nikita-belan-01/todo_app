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
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return zero, domain.NewBadRequestError("unexpected content after json body", err)
	}

	return dto.BuildDomain(), nil
}

func decodeError(err error) error {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return domain.NewBadRequestError(fmt.Sprintf("request body must not exceed %d bytes", maxBytesErr.Limit), err)
	}

	if errors.Is(err, io.EOF) {
		return domain.NewBadRequestError("empty request body", err)
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		return domain.NewBadRequestError("malformed json body: unexpected end of input", err)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return domain.NewBadRequestError(
			fmt.Sprintf("field %q: expected %s, got %s", typeErr.Field, typeErr.Type, typeErr.Value), err)
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return domain.NewBadRequestError(
			fmt.Sprintf("malformed json body at offset %d", syntaxErr.Offset), err)
	}

	if strings.HasPrefix(err.Error(), unknownFieldPrefix) {
		return domain.NewBadRequestError(
			fmt.Sprintf("unknown field %s", strings.Trim(err.Error(), unknownFieldPrefix)), err)
	}

	return domain.NewBadRequestError("invalid json body request", err)
}
