package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL string
	Port   string
	Secret string
}

func Load() (*Config, error) {
	// 1. Attempt to load .env file
	// We don't log.Fatal here because in production (Railway),
	// variables are injected directly and a .env file won't exist.
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323" // Default fallback
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		return nil, fmt.Errorf("SECRET is not set")
	}

	return &Config{
		DB_URL: dbURL,
		Port:   port,
		Secret: secret,
	}, nil
}

func ConfigLogger() {
	logfile := os.Getenv("LOGFILE")
	var output io.Writer

	if logfile == "" {
		// NO LOGFILE env set, default to STDOUT (not running locally)
		output = os.Stdout
		fmt.Println("config.go->ConfigLogger(): No LOGFILE env set. Logging to STDOUT.")
	} else {
		// Logfile is set, attempt to open it (Local behavior)
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("Error opening logfile %s: %v. Falling back to STDOUT.\n", logfile, err)
			output = os.Stdout
		} else {
			output = file
			fmt.Println("config.go->ConfigLogger(): Logging to file:", logfile)
		}
	}

	// Create the JSON handler using the determined output
	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("Logger initialized successfully")
}
