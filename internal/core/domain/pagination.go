package domain

import (
	"errors"
	"fmt"
)

var ErrInvalidPagination = errors.New("invalid pagination")

type Pagination struct {
	Page  int
	Limit int
}

type PaginationError struct {
	field    string
	max, got int
}

func (e *PaginationError) Error() string {
	return fmt.Sprintf("must be between 1 and %d, got %d", e.max, e.got)
}

func (e *PaginationError) Unwrap() error {
	return ErrInvalidPagination
}

func (e *PaginationError) Field() string {
	return e.field
}

func (p Pagination) validate() []*PaginationError {
	var errs []*PaginationError

	if p.Page < 1 || p.Page > PageMaxValue {
		errs = append(errs, &PaginationError{
			field: "page",
			max:   PageMaxValue,
			got:   p.Page})
	}

	if p.Limit < 1 || p.Limit > LimitMaxValue {
		errs = append(errs, &PaginationError{
			field: "limit",
			max:   LimitMaxValue,
			got:   p.Limit})
	}

	return errs
}

func (p Pagination) Validate() error {
	var v Validation

	for _, err := range p.validate() {
		v.Add(err.field, err)
	}

	return v.Err()
}
