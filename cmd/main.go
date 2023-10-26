package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"analytics_collector/internal/api/http-server/handlers/metrics/analytics"
	metricWorkerPool "analytics_collector/internal/api/http-server/handlers/metrics/worker_pool"
	"analytics_collector/internal/config"
	sl "analytics_collector/internal/logging"
	"analytics_collector/internal/storage/postgres"
)

func main() {
	// create context
	appCtx := context.Background()
	appCtx, cancel := signal.NotifyContext(appCtx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// parse config
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config is not found. Error: %s", err)
	}

	// create storage
	storage, err := postgres.New(appCtx, cfg.DB)
	if err != nil {
		log.Fatalf("storage is not created. Error: %s", err)
	}

	// create logging
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("logging is not created. Error: %s", err)
	}

	mux := http.NewServeMux()

	// create workers for request "/analytics"
	jobsChannel := metricWorkerPool.New(appCtx, logger, cfg.Server.WorkersCount, storage)
	mux.HandleFunc("/analytics", analytics.HandleAnalytics(appCtx, logger, jobsChannel))

	server := &http.Server{
		Addr:    cfg.GetServerAddr(),
		Handler: mux,
	}

	// start server
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server not started. Error: %s", err)
		}
	}()

	logger.Info("Server started",
		slog.String("env", cfg.Env),
		slog.String("Address", cfg.GetServerAddr()),
	)

	<-appCtx.Done()
	logger.Info("Service stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed. Error: %+v", err)
	}

	<-ctx.Done()

	logger.Info("Program correctly finished")
}
