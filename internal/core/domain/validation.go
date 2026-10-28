package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nikita-belan-01/todo_app/pkg/validator"
)

type Validation struct {
	problems []string
	causes   []error
}

func (v *Validation) addProblem(field, msg string) {
	if field == "" {
		v.problems = append(v.problems, msg)
		return
	}

	v.problems = append(v.problems, fmt.Sprintf("%s: %s", field, msg))
}

func (v *Validation) Add(field string, err error) {
	if err == nil {
		return
	}

	for line := range strings.SplitSeq(err.Error(), "\n") {
		if line == "" {
			continue
		}

		v.addProblem(field, line)
	}
}

func (v *Validation) Addc(field string, cause error, format string, args ...any) {
	if cause != nil {
		v.causes = append(v.causes, cause)
	}

	v.addProblem(field, fmt.Sprintf(format, args...))
}

func (v *Validation) Addf(field, format string, args ...any) {
	v.addProblem(field, fmt.Sprintf(format, args...))
}

func (v *Validation) AddNullableBool(field string, n Nullable[bool]) {
	if n.Set && n.Value == nil {
		v.Addc(field, ErrNullNotAllowed, "must not be null")
	}
}

func (v *Validation) AddNullableLen(field string, n Nullable[string], min, max int) {
	if !n.Set {
		return
	}

	if n.Value == nil {
		v.Addc(field, ErrNullNotAllowed, "must not be null")
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
		v.Addc(field, ErrNullNotAllowed, "must not be null")
		return
	}

	v.Add(field, validate(*n.Value))
}

func (v *Validation) Err() error {
	if len(v.problems) == 0 {
		return nil
	}

	causes := append([]error{ErrInvalidArgument}, v.causes...)

	return NewAppError(KindInvalidArgument, strings.Join(v.problems, "; "), errors.Join(causes...))
}
