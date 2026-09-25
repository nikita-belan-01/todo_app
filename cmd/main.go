package main

import (
	"context"
	"flag"
	"log"
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

	run(configPath)
}

func run(confPath string) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf, err := config.New(confPath)
	if err != nil {
		log.Fatalf("init config: %v", err)
	}

	logger, err := logger.New(conf.Logger.Dir, conf.Logger.Level)
	if err != nil {
		log.Fatalf("init app logger: %v", err)
	}
	defer func() {
		if err := logger.Close(); err != nil {
			log.Fatalf("logger close: %v", err)
		}
	}()

	logger.Debug("initializing postgres connection pool")

	poolCtx, poolCancel := context.WithTimeout(ctx, conf.Postgres.PingTimeout)
	pool, err := pool.NewConnectionPool(poolCtx, &conf.Postgres)
	if err != nil {
		logger.Fatal("init postgres pool", zap.Error(err))
	}
	defer pool.Close()
	poolCancel()

	logger.Debug("initializing features", zap.String("feature", "users"))

	usersRepository := users_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersDeliveryHTTP := users_delivery_http.NewUsersHTTPHandler(&conf.Handler, usersService)

	logger.Debug("initializing HTTP server")

	httpServer := server.New(&conf.HTTPServer, logger,
		middleware.RequestID(),
		middleware.Logger(logger.With(zap.String("component", "middleware"))),
		middleware.Panic(),
		middleware.Trace())

	apiVersionRouter := server.NewApiVersionRouter(server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersDeliveryHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
