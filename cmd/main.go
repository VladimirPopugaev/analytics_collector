package main

import (
	"log"

	"analytics_collector/internal/config"
	sl "analytics_collector/internal/logger"
)

const (
	configPath = "./configs/local.yaml"
)

func main() {
	conf, err := config.New(configPath)
	if err != nil {
		log.Fatalf("Config is not found. Error: %s", err)
	}
	log.Printf("config: = %v", conf)

	logger, err := sl.SetupLogger(conf.Env)
	if err != nil {
		log.Fatalf("Logger is not created. Error: %s", err)
	}
	_ = logger

	//TODO: start service
}
