package config

import (
	"fmt"

	"github.com/nikita-belan-01/todo_app/pkg/default_env"
)

type Config struct {
	Logger     Logger
	HTTPServer HTTPServer
	Handler    Handler
	Postgres   Postgres
	LocalTime  string
}

func New(confPath string) (*Config, error) {
	if err := loadEnv(confPath); err != nil {
		return nil, fmt.Errorf("load env by path: %s, %w", confPath, err)
	}

	logger, err := newLogger()
	if err != nil {
		return nil, err
	}

	httpServer, err := newHTTPServer()
	if err != nil {
		return nil, err
	}

	handler, err := newHandler()
	if err != nil {
		return nil, err
	}

	postgres, err := newPostgres()
	if err != nil {
		return nil, err
	}

	localTime := default_env.GetString("LOCAL_TIME", "UTC")

	return &Config{
		Logger:     logger,
		HTTPServer: httpServer,
		Handler:    handler,
		Postgres:   postgres,
		LocalTime:  localTime,
	}, nil
}
