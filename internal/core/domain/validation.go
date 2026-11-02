package domain

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/nikita-belan-01/todo_app/pkg/validator"
)

const maxProblems = 10

type Problem struct {
	Field   string
	Message string
}

func (p Problem) String() string {
	if p.Field == "" {
		return p.Message
	}

	return fmt.Sprintf("%s: %s", p.Field, p.Message)
}

type Validation struct {
	problems []Problem
	causes   []error
	skipped  int
}

func (v *Validation) addProblem(field, msg string) {
	if len(v.problems) >= maxProblems {
		v.skipped++

		return
	}

	v.problems = append(v.problems, Problem{Field: field, Message: msg})
}

func (v Validation) Problems() []Problem {
	return v.problems
}

func (v *Validation) Add(field string, err error) {
	if err == nil {
		return
	}

	v.addCause(err)

	text := err.Error()

	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Message() != "" {
		text = appErr.Message()
	}

	for line := range strings.SplitSeq(text, "\n") {
		if line == "" {
			continue
		}

		v.addProblem(field, line)
	}
}

func (v *Validation) AddAll(err error) {
	v.Add("", err)
}

func (v *Validation) AddWithCause(field string, cause error, format string, args ...any) {
	v.addCause(cause)
	v.addProblem(field, fmt.Sprintf(format, args...))
}

func (v *Validation) addCause(cause error) {
	if cause == nil || slices.Contains(v.causes, cause) {
		return
	}

	v.causes = append(v.causes, cause)
}

func (v *Validation) AddNullableBool(field string, n Nullable[bool]) {
	if n.Set && n.Value == nil {
		v.AddWithCause(field, ErrNullNotAllowed, "must not be null")
	}
}

func (v *Validation) AddNullableLen(field string, n Nullable[string], min, max int) {
	if !n.Set {
		return
	}

	if n.Value == nil {
		v.AddWithCause(field, ErrNullNotAllowed, "must not be null")
		return
	}

	v.Add(field, validator.ValidateLen(*n.Value, min, max))
}

func (v *Validation) AddOptionalLen(field string, n Nullable[string], min, max int) {
	if !n.Set || n.Value == nil {
		return
	}

	v.Add(field, validator.ValidateLen(*n.Value, min, max))
}

func (v *Validation) AddNullableFunc(field string, n Nullable[string], validate func(string) error) {
	if !n.Set {
		return
	}

	if n.Value == nil {
		v.AddWithCause(field, ErrNullNotAllowed, "must not be null")
		return
	}

	v.Add(field, validate(*n.Value))
}

func (v *Validation) Err() error {
	if len(v.problems) == 0 {
		return nil
	}

	parts := make([]string, 0, len(v.problems))
	for _, p := range v.problems {
		parts = append(parts, p.String())
	}

	message := strings.Join(parts, "; ")
	if v.skipped > 0 {
		message = fmt.Sprintf("%s; and %d more problems", message, v.skipped)
	}

	return NewValidationError(message, v.problems, errors.Join(v.causes...))
}
