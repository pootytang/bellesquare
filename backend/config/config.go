package config

import (
	"fmt"
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
	if logfile == "" {
		logfile = "app.log" // Default log file name
	}
	fmt.Println("config.go->ConfigLogger(): Logfile set to:", logfile)

	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
	} else {
		// Create a JSON handler that writes to the file
		handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true, // Include source file and line number
		})

		logger := slog.New(handler)

		// Set as the default logger
		slog.SetDefault(logger)
		slog.Info("Logger initialized successfully")
	}
}
