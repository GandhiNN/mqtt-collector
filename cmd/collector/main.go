package main

import (
	"log"
	"mqtt-collector/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Brokers) == 0 {
		log.Fatalf("No brokers configured")
	}

	log.Printf("Starting MQTT Topic Collector")
	log.Printf("Database Service: %s", cfg.DBServiceURL)
	log.Printf("Duration: %v", cfg.CollectionDuration)
	log.Printf("Configured brokers: %d", len(cfg.Brokers))

	mc := collector.NewMultiCollector(cfg.Brokers, cfg.DBServiceURL)

	if err := mc.Run(cfg.CollectionDuration); err != nil {
		log.Fatalf("Collector error: %v", err)
	}
}
