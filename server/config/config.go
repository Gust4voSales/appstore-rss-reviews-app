package config

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		// TODO add polling interval
	}

	log.Printf("📦 Config loaded. PORT=%s", port) // PRINTING CONFIG FOR DEBUGGING PURPOSES, WOULDN'T LOG SENSITIVE DATA IN PRODUCTION ON REAL APP
	return &Config{
		Port: port,
	}
}
