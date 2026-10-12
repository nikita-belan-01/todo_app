package config

import "fmt"

type Config struct {
	Logger     Logger
	HTTPServer HTTPServer
	Handler    Handler
	Postgres   Postgres
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

	return &Config{
		Logger:     logger,
		HTTPServer: httpServer,
		Handler:    handler,
		Postgres:   postgres,
	}, nil
}
