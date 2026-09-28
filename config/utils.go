package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

var (
	ErrEmptySecret     = errors.New("empty secret")
	ErrInvalidPort     = errors.New("parameter must be between 1 and 65535")
	ErrInvalidDuration = errors.New("invalid duration")
)

func loadEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return nil
}

func getSecretFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	secret := strings.TrimRight(string(data), "\r\n")
	if secret == "" {
		return "", ErrEmptySecret
	}

	return secret, nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return ErrInvalidPort
	}

	return nil
}

func validateDuration(duration time.Duration) error {
	if duration <= 0 {
		return ErrInvalidDuration
	}

	return nil
}
