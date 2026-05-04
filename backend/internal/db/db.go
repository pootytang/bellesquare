package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect establishes a connection to the PostgreSQL database using the provided URL and returns a connection pool.
// Pooling allows for efficient management of database connections, enabling multiple concurrent queries without the overhead of establishing a new connection for each query.
// In short, multiple connections are available in the pool. When a query is executed, a connection from the pool is used, and once the query is complete, the connection is returned to the pool for reuse.
func Connect(dburl string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(dburl)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to parse database url: %v", err))
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to create connection pool: %v", err))
		return nil, err
	}

	// Test the connection
	err = pool.Ping(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to connect to database: %v", err))
		pool.Close()
		return nil, err
	}

	slog.Info("Successfully connected to database")
	return pool, nil
}