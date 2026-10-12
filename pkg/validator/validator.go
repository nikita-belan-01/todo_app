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

func ValidateLen(s string, min, max int) error {
	rs := []rune(s)
	if len(rs) < min || len(rs) > max {
		return fmt.Errorf("%w: len must be between %d and %d", ErrInvalidLen, min, max)
	}

	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	var errs []error
	if err := ValidateLen(phoneNumber, 10, 15); err != nil {
		errs = append(errs, fmt.Errorf("phone number: %w", err))
	}

	re := regexp.MustCompile(`^\+[0-9]+$`)

	if !re.MatchString(phoneNumber) {
		errs = append(errs, fmt.Errorf("phone number: invalid format: %w", ErrInvalidArgument))
	}

	return errors.Join(errs...)
}
