package domain

import "github.com/nikita-belan-01/todo_app/pkg/validator"

type Nullable[T any] struct {
	Value *T
	Set   bool
}

func applyNotNull[T comparable](dst *T, n Nullable[T]) bool {
	if !n.Set || n.Value == nil || *n.Value == *dst {
		return false
	}

	*dst = *n.Value

	return true
}

func applyNullable[T comparable](dst **T, n Nullable[T]) bool {
	if !n.Set || equalPtr(n.Value, *dst) {
		return false
	}

	*dst = n.Value

	return true
}

func equalPtr[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b
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
