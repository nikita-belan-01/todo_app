package http

import (
	"context"
	"net/http"

	"github.com/nikita-belan-01/todo_app/config"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/server"
)

type HealthHTTPHandler struct {
	Config *config.Handler
	Checks []Check
}

type Check struct {
	Name   string
	Pinger Pinger
}

type Pinger interface {
	Ping(ctx context.Context) error
}

func NewHTTPHealthHandler(config *config.Handler, databasePinger Pinger) *HealthHTTPHandler {
	return &HealthHTTPHandler{
		Config: config,
		Checks: []Check{{
			Name:   "database",
			Pinger: databasePinger,
		}}}
}

func (h *HealthHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodGet, Path: "/health", Handler: h.Health},
		{Method: http.MethodGet, Path: "/ready", Handler: h.Ready},
	}
}
