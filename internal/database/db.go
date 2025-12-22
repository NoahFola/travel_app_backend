package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB creates a connection pool
func InitDB() *pgxpool.Pool {
	// Example DSN: "postgres://user:password@localhost:5432/dbname"
	dsn := os.Getenv("DATABASE_URL")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	// Optimization: Max connections (adjust based on your DB size)
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	// Create the pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Ping to verify connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return pool
}
