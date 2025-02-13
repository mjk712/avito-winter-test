package main

import (
	"avito-winter-test/internal/config"
	httpServer "avito-winter-test/internal/http-server"
	"avito-winter-test/internal/service"
	"avito-winter-test/internal/storage"
	"avito-winter-test/internal/tools"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	//config
	cfg := config.New()
	//log
	log := setupLogger(cfg.Env)
	log.Info(
		"starting server",
		slog.String("env", cfg.Env),
		slog.String("version", "1"),
	)
	log.Debug("debug messages are enabled")
	//ctx
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//storage
	repo, err := storage.New(cfg.ConnectionString)
	if err != nil {
		log.Error("failed to init storage", tools.ErrAttr(err))
		os.Exit(1)
	}
	//services
	merchShopService := service.NewMerchShopService(repo)

	//server
	serv := httpServer.NewServer(mainCtx, log, cfg, merchShopService)

	log.Info("starting server")

	//Graceful shutdown

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

		os.Exit(1)
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
