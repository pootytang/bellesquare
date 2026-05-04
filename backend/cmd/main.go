package main

import (
	"bellesquare-be/config"
	"bellesquare-be/internal/db"
	"bellesquare-be/internal/routes"
	"bellesquare-be/internal/ws"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	config.ConfigLogger()

	pool, err := db.Connect(cfg.DB_URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	e := echo.New()
	// 1. Initialize and start the websocket Hub
	hub := ws.NewHub()
	go hub.Run()

	routes.SetupRoutes(e, pool, cfg.Secret, hub)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: e,
	}

	// Start server in a goroutine
	go func() {
		e.Logger.Info("Starting server on port " + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("shuting down the server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // Railway sends SIGTERM
	<-quit

	e.Logger.Info("Shutting down server...")

	// 4. Standard library context-based shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		e.Logger.Error("Server forced to shutdown", "error", err)
	}

	e.Logger.Info("Server exited gracefully")
}
