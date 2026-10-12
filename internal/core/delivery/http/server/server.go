package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/nikita-belan-01/todo_app/config"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/middleware"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
	"go.uber.org/zap"
)

type Server struct {
	mux    *http.ServeMux
	conf   *config.HTTPServer
	logger *logger.Logger

	middleware []middleware.Middleware
}

func New(
	conf *config.HTTPServer,
	logger *logger.Logger,
	middleware ...middleware.Middleware) *Server {
	return &Server{
		mux:        http.NewServeMux(),
		conf:       conf,
		logger:     logger.With(zap.String("component", "server")),
		middleware: middleware,
	}
}

func (s Server) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(s.mux, s.middleware...)
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", s.conf.Port),
		ReadTimeout:    s.conf.ReadTimeout,
		WriteTimeout:   s.conf.WriteTimeout,
		IdleTimeout:    s.conf.IdleTimeout,
		MaxHeaderBytes: s.conf.MaxHeaderBytes,
		Handler:        mux,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("listen: %s: %w", server.Addr, err)
	}

	s.logger.Debug("HTTP server listening", zap.String("addr", ln.Addr().String()))

	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		if err := server.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
		s.logger.Debug("shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.conf.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}

		s.logger.Debug("HTTP server stopped")

		return nil
	}
}

func (s *Server) RegisterAPIRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := fmt.Sprintf("/api/%s", router.apiVersion)

		s.mux.Handle(
			fmt.Sprintf("%s/", prefix),
			http.StripPrefix(prefix, router),
		)
	}
}
