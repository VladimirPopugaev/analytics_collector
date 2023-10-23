package main

import (
	"analytics_collector/internal/config"
	sl "analytics_collector/internal/logger"
	"analytics_collector/internal/storage/postgres"
	"context"
	syslog "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	configPath = "./configs/local.yaml"
)

func main() {
	// create context
	appCtx := context.Background()
	appCtx, cancel := signal.NotifyContext(appCtx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// parse config
	cfg, err := config.New(configPath)
	if err != nil {
		syslog.Fatalf("Config is not found. Error: %s", err)
	}

	// create storage
	storage, err := postgres.New(appCtx, cfg.DB)
	if err != nil {
		syslog.Fatalf("storage is not created. Error: %s", err)
	}
	syslog.Printf("storage started")
	_ = storage

	log, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		syslog.Fatalf("logger is not created. Error: %s", err)
	}

	//TODO: add handlers
	server := &http.Server{
		Addr:    cfg.GetServerAddr(),
		Handler: nil,
	}

	log.Info("Server started",
		slog.String("env", cfg.Env),
		slog.String("Address", cfg.GetServerAddr()),
	)

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start service")
	}

	log.Info("Service stopped")
}
