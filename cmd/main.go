package main

import (
	"context"
	syslog "log"
	"log/slog"
	"net/http"
	"time"

	"analytics_collector/internal/config"
	sl "analytics_collector/internal/logger"
	"analytics_collector/internal/storage/postgres"
)

const (
	configPath = "./configs/local.yaml"
)

func main() {
	cfg, err := config.New(configPath)
	if err != nil {
		syslog.Fatalf("Config is not found. Error: %s", err)
	}

	storage, err := postgres.New(cfg.DB)
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

	appContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = appContext

	log.Info("Server started",
		slog.String("env", cfg.Env),
		slog.String("Address", cfg.GetServerAddr()),
	)

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start service")
	}

	log.Info("Service stopped")
}
