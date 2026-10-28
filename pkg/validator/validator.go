package validator

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrInvalidLen      = errors.New("invalid len")
	ErrInvalidArgument = errors.New("invalid argument")
)

type LenError struct {
	Min, Max, Got int
}

func (e *LenError) Error() string {
	return fmt.Sprintf("length must be between %d and %d characters, got %d", e.Min, e.Max, e.Got)
}

func (e *LenError) Unwrap() error {
	return ErrInvalidLen
}

type FormatError struct {
	Want string
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("invalid format, expected %s", e.Want)
}

func (e *FormatError) Unwrap() error {
	return ErrInvalidArgument
}

var phoneNumberRe = regexp.MustCompile(`^\+[0-9]+$`)

func ValidateLen(s string, min, max int) error {
	if n := len([]rune(s)); n < min || n > max {
		return &LenError{Min: min, Max: max, Got: n}
	}

	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	var errs []error
	if err := ValidateLen(phoneNumber, 10, 15); err != nil {
		errs = append(errs, err)
	}

	if !phoneNumberRe.MatchString(phoneNumber) {
		errs = append(errs, &FormatError{Want: "a leading + followed by digits only"})
	}

	return errors.Join(errs...)
}
