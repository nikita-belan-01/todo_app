package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikita-belan-01/todo_app/config"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/middleware"
	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/server"
	"github.com/nikita-belan-01/todo_app/internal/core/repository/postgres/pool"
	users_delivery_http "github.com/nikita-belan-01/todo_app/internal/features/users/delivery/http"
	users_postgres "github.com/nikita-belan-01/todo_app/internal/features/users/repository/postgres"
	users_service "github.com/nikita-belan-01/todo_app/internal/features/users/service"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", ".env", "config file path")
	flag.StringVar(&configPath, "c", ".env", "config file path")

	flag.Parse()

	if err := run(configPath); err != nil {
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}
}

func run(configPath string) error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf, err := config.New(configPath)
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	appLogger, err := logger.New(conf.Logger.Dir, conf.Logger.Level)
	if err != nil {
		return fmt.Errorf("init app logger: %w", err)
	}
	defer func() {
		if err := appLogger.Close(); err != nil {
			log.Printf("close app logger: %v", err)
		}
	}()

	appLogger.Debug("initializing postgres connection pool")

	poolCtx, poolCancel := context.WithTimeout(ctx, conf.Postgres.PingTimeout)
	defer poolCancel()

	pool, err := pool.NewConnectionPool(poolCtx, &conf.Postgres)
	if err != nil {
		return fmt.Errorf("init postgres pool: %w", err)
	}
	defer pool.Close()

	appLogger.Debug("initializing features", zap.String("feature", "users"))

	usersRepository := users_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersDeliveryHTTP := users_delivery_http.NewUsersHTTPHandler(&conf.Handler, usersService)

	appLogger.Debug("initializing HTTP server")

	httpServer := server.New(&conf.HTTPServer, appLogger,
		middleware.RequestID(),
		middleware.Logger(appLogger.With(zap.String("component", "middleware"))),
		middleware.Panic(),
		middleware.BodyLimit(conf.HTTPServer.MaxBodyBytes),
		middleware.Trace())

	apiVersionRouter := server.NewApiVersionRouter(server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersDeliveryHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}
