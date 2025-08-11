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
	StorageFilePath string
}

func Load() *Config {
	port := os.Getenv("PORT")
	pollingIntervalSecondsStr := os.Getenv("POLLING_INTERVAL_SECONDS")
	storageFilePath := os.Getenv("STORAGE_FILE_PATH")
	appID := os.Getenv("APP_ID")

	if port == "" {
		port = "8080"
	}

	if pollingIntervalSecondsStr == "" {
		pollingIntervalSecondsStr = "30"
	}

	if appID == "" {
		appID = "447188370" // Default to Snapchat app ID
	}

	if storageFilePath == "" {
		storageFilePath = "data/reviews-" + appID + ".json"
	}

	pollingIntervalSeconds, err := strconv.Atoi(pollingIntervalSecondsStr)
	if err != nil {
		log.Fatalf("invalid polling interval seconds: %v", err)
	}

	// PRINTING CONFIG FOR DEBUGGING PURPOSES, WOULDN'T LOG SENSITIVE DATA IN PRODUCTION ON REAL APP
	log.Printf("📦 Config loaded. PORT=%s, POLLING_INTERVAL_SECONDS=%d, APP_ID=%s, STORAGE_FILE_PATH=%s", port, pollingIntervalSeconds, appID, storageFilePath)

	return &Config{
		Port:            port,
		PollingInterval: time.Duration(pollingIntervalSeconds) * time.Second,
		AppID:           appID,
		StorageFilePath: storageFilePath,
	}
}
