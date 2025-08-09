package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	PollingInterval time.Duration
	AppID           string
}

func Load() *Config {
	// appID := "835599320" // Hardcoded, later add dynamic app id?
	appID := "389801252" // Hardcoded, later add dynamic app id?
	// appID := "1481979331" // Hardcoded, later add dynamic app id?
	// appID := "447188370" // Hardcoded, later add dynamic app id?

	port := os.Getenv("PORT")
	pollingIntervalSecondsStr := os.Getenv("POLLING_INTERVAL_SECONDS")

	if port == "" {
		port = "8080"
	}

	if pollingIntervalSecondsStr == "" {
		pollingIntervalSecondsStr = "15"
	}

	pollingIntervalSeconds, err := strconv.Atoi(pollingIntervalSecondsStr)
	if err != nil {
		log.Fatalf("invalid polling interval seconds: %v", err)
	}

	// PRINTING CONFIG FOR DEBUGGING PURPOSES, WOULDN'T LOG SENSITIVE DATA IN PRODUCTION ON REAL APP
	log.Printf("📦 Config loaded. PORT=%s, POLLING_INTERVAL_SECONDS=%d, APP_ID=%s", port, pollingIntervalSeconds, appID)

	return &Config{
		Port:            port,
		PollingInterval: time.Duration(pollingIntervalSeconds) * time.Second,
		AppID:           appID,
	}
}
