package main

import (
	"context"
	"delimed/internal/config"
	"delimed/internal/repository/postgres"
	"delimed/internal/service"
	"delimed/internal/transport/handler"
	"delimed/internal/transport/httpserver"
	sl "delimed/internal/utils/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Раскомментировать после генерации swagger документации:
	_ "delimed/docs" // swagger docs
)

func main() {
	config := config.MustLoadConfig()
	logger := sl.SetupLogger(config.Env)

	logger.Info("ockc", slog.String("env", config.Env))
	logger.Debug("debug message are enabled")

	db, err := postgres.NewStorage(config.Database.DSN,
		config.Database.ConnMaxIdleTime,
		config.Database.ConnMaxLifetime,
		config.Database.MaxOpenConns,
		config.Database.MaxIdleConns,
	)
	if err != nil {
		logger.Error("Failed to connect to database", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("database connect success")

	service := service.NewService(
		db,
		config.JWTSecret,
		logger,
		config.CDEK.ClientID,
		config.CDEK.ClientSecret,
		config.Dellin.AppKey,
	)

	logger.Info("service init success")

	handler := handler.NewHandler(service)

	logger.Info("hanlder init success")

	server := httpserver.NewServer(config.JWTSecret,
		handler,
		config.HTTPServer.Port,
		config.HTTPServer.WriteTimeout,
		config.HTTPServer.ReadTimeout,
		config.HTTPServer.IdleTimeout)

	logger.Info("server config init success",
		slog.String("Addr", config.HTTPServer.Port),
		slog.String("WriteTimeout", config.HTTPServer.WriteTimeout.String()),
		slog.String("ReadTimeout", config.HTTPServer.ReadTimeout.String()),
		slog.String("IdleTimeout", config.HTTPServer.IdleTimeout.String()),
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", sl.Err(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", sl.Err(err))
	}

	logger.Info("Server exited")
}
