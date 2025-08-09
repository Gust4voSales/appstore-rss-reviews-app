package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/config"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/api"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/app"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/crons/appstore_reviews_poller"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/repositories"
)

func main() {
	// Load config
	cfg := config.Load()

	// Load repositories
	repositories := repositories.Load()

	// Load app service
	appService := app.New(repositories, cfg)

	// Load cron jobs
	poller := appstore_reviews_poller.New(cfg, appService)
	// Run cron jobs
	go poller.Run(context.Background())

	// Load Gin router setup
	router := api.NewRouter(appService)

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Run server in goroutine
	go func() {
		log.Printf("🚀 Server running on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
