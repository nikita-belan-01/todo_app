package config

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/nikita-belan-01/todo_app/pkg/default_env"
)

var (
	ErrEmptyPostgresUser         = errors.New("empty POSTGRES_USER")
	ErrEmptyPostgresDB           = errors.New("empty POSTGRES_DB")
	ErrEmptyPostgresPasswordFile = errors.New("empty POSTGRES_PASSWORD_FILE")
)

type Postgres struct {
	Host        string
	Port        int
	Username    string
	Password    string
	Database    string
	PingTimeout time.Duration
}

func newPostgres() (Postgres, error) {
	host := default_env.GetString("POSTGRES_HOST", "127.0.0.1")

	port, err := default_env.GetNumber("POSTGRES_PORT", 5432)
	if err != nil {
		return Postgres{}, err
	}

	if err = validatePort(port); err != nil {
		return Postgres{}, fmt.Errorf("POSTGRES_PORT: %w", err)
	}

	username := default_env.GetString("POSTGRES_USER", "")
	if username == "" {
		return Postgres{}, ErrEmptyPostgresUser
	}

	var password string
	passwordFile := default_env.GetString("POSTGRES_PASSWORD_FILE", "")
	if passwordFile == "" {
		return Postgres{}, ErrEmptyPostgresPasswordFile
	}

	password, err = getSecretFromFile(passwordFile)
	if err != nil {
		return Postgres{}, fmt.Errorf("Postgres password: %w", err)
	}

	database := default_env.GetString("POSTGRES_DB", "")
	if database == "" {
		return Postgres{}, ErrEmptyPostgresDB
	}

	pingTimeout, err := default_env.GetDuration("POSTGRES_PING_TIMEOUT", 5*time.Second)
	if err != nil {
		return Postgres{}, err
	}

	if err = validateDuration(pingTimeout); err != nil {
		return Postgres{}, fmt.Errorf("POSTGRES_PING_TIMEOUT: %w", err)
	}

	return Postgres{
		Host:        host,
		Port:        port,
		Username:    username,
		Password:    password,
		Database:    database,
		PingTimeout: pingTimeout,
	}, nil
}

func (p Postgres) ToDSN() string {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(p.Username, p.Password),
		Host:     fmt.Sprintf("%s:%d", p.Host, p.Port),
		Path:     "/" + p.Database,
		RawQuery: "sslmode=disable",
	}

	return dsn.String()
}
