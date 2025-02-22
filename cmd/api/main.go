package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avito-winter-test/internal/config"
	httpServer "avito-winter-test/internal/http-server"
	"avito-winter-test/internal/service"
	"avito-winter-test/internal/storage"
	"avito-winter-test/internal/tools"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	// Config
	cfg := config.New()
	// Log
	log := setupLogger(cfg.Env)
	log.Info(
		"starting server",
		slog.String("env", cfg.Env),
		slog.String("version", "1"),
	)
	log.Debug("debug messages are enabled")
	// Context
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Storage
	repo, err := storage.New(cfg.ConnectionString)
	if err != nil {
		log.Error("failed to init storage", tools.ErrAttr(err))
		return
	}
	// Services
	services := &httpServer.Services{
		MerchShop: service.NewMerchShopService(repo),
		Auth:      service.NewAuthService(repo),
	}

	// Server
	serv := httpServer.NewServer(mainCtx, log, cfg, services)

	log.Info("starting server")

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := serv.ListenAndServe(); err != nil {
			log.Error("failed to start server", tools.ErrAttr(err))
		}
	}()
	log.Info("server started")

	<-stop
	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", tools.ErrAttr(err))
		return
	}
	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
