package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Galiwer/grapi/internal/handlers"
)

func main() {
	version := getEnv("APP_VERSION", "1.1.0")
	environment := getEnv("APP_ENV", "local")

	mux := http.NewServeMux()

	api := handlers.New(version, environment)

	mux.HandleFunc("GET /", api.Home)
	mux.HandleFunc("GET /health", api.Health)
	mux.HandleFunc("GET /ready", api.Ready)
	mux.HandleFunc("GET /version", api.Version)
	mux.HandleFunc("GET /api/message", api.Message)
	mux.HandleFunc("GET /metrics", api.Metrics)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("DevOps Demo API running on port 8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
