package config

import (
	"fmt"
	"time"

	"github.com/nikita-belan-01/todo_app/pkg/default_env"
)

const (
	defaultMaxHeaderBytes int = 5 << 20
)

type HTTPServer struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	MaxHeaderBytes  int
}

func newHTTPServer() (HTTPServer, error) {
	port, err := default_env.GetNumber("HTTP_SERVER_PORT", 8080)
	if err != nil {
		return HTTPServer{}, err
	}

	if err = validatePort(port); err != nil {
		return HTTPServer{}, fmt.Errorf("HTTP_SERVER_PORT: %w", err)
	}

	readTimeout, err := default_env.GetDuration("HTTP_SERVER_READ_TIMEOUT", 15*time.Second)
	if err != nil {
		return HTTPServer{}, err
	}

	if err = validateDuration(readTimeout); err != nil {
		return HTTPServer{}, fmt.Errorf("HTTP_SERVER_READ_TIMEOUT: %w", err)
	}

	idleTimeout, err := default_env.GetDuration("HTTP_SERVER_IDLE_TIMEOUT", time.Minute)
	if err != nil {
		return HTTPServer{}, err
	}

	if err = validateDuration(idleTimeout); err != nil {
		return HTTPServer{}, fmt.Errorf("HTTP_SERVER_IDLE_TIMEOUT: %w", err)
	}

	writeTimeout, err := default_env.GetDuration("HTTP_SERVER_WRITE_TIMEOUT", 15*time.Second)
	if err != nil {
		return HTTPServer{}, err
	}

	if err = validateDuration(writeTimeout); err != nil {
		return HTTPServer{}, fmt.Errorf("HTTP_WRITE_READ_TIMEOUT: %w", err)
	}

	shutdownTimeout, err := default_env.GetDuration("HTTP_SERVER_SHUTDOWN_TIMEOUT", 5*time.Second)
	if err != nil {
		return HTTPServer{}, err
	}

	if err = validateDuration(shutdownTimeout); err != nil {
		return HTTPServer{}, fmt.Errorf("HTTP_SERVER_SHUTDOWN_TIMEOUT: %w", err)
	}

	maxHeaderBytes, err := default_env.GetNumber("HTTP_SERVER_MAX_HEADER_BYTES", defaultMaxHeaderBytes)
	if err != nil {
		return HTTPServer{}, err
	}

	return HTTPServer{
		Port:            port,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		ShutdownTimeout: shutdownTimeout,
		MaxHeaderBytes:  maxHeaderBytes,
		IdleTimeout:     idleTimeout,
	}, nil
}
