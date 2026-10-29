package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nikita-belan-01/todo_app/internal/core/domain"
	"github.com/nikita-belan-01/todo_app/pkg/strutil"
)

const MaxEchoedValue = 64

const unknownFieldPrefix = "json: unknown field"

func DecodeOptional[T interface{ BuildDomain() (DT, error) }, DT any](r *http.Request) (DT, error) {
	return decode[T](r, true)
}

func Decode[T interface{ BuildDomain() (DT, error) }, DT any](r *http.Request) (DT, error) {
	return decode[T](r, false)
}

func decode[T interface{ BuildDomain() (DT, error) }, DT any](r *http.Request, allowEmpty bool) (DT, error) {
	var dto *T
	var zero DT

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dto); err != nil {
		if allowEmpty && errors.Is(err, io.EOF) {
			var empty T

			return empty.BuildDomain()
		}

		return zero, decodeError(err)
	}

	if dto == nil {
		return zero, domain.NewAppError(domain.KindInvalidArgument, "request body must be a json object", nil)
	}

	var extra json.RawMessage
	switch err := dec.Decode(&extra); {
	case err == nil:
		return zero, domain.NewAppError(domain.KindInvalidArgument, "unexpected content after json body", nil)
	case !errors.Is(err, io.EOF):
		return zero, decodeError(err)
	}

	return (*dto).BuildDomain()
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
		if typeErr.Field == "" {
			return domain.NewAppError(
				domain.KindInvalidArgument,
				fmt.Sprintf("request body must be a json object, got %s", strutil.Shorten(typeErr.Value, MaxEchoedValue)),
				err)
		}
		return domain.NewAppError(
			domain.KindInvalidArgument,
			fmt.Sprintf("field %q: unexpected %s", typeErr.Field, strutil.Shorten(typeErr.Value, MaxEchoedValue)), err)
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
			fmt.Sprintf("unknown field %s", strutil.Shorten(strings.TrimSpace(field), MaxEchoedValue)), err)
	}

	return domain.NewAppError(domain.KindInvalidArgument,
		fmt.Sprintf("malformed value in request body: %s", strutil.Shorten(err.Error(), MaxEchoedValue)), err)
}
