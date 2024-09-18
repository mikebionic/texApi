package database

import (
	"context"
	"fmt"
	"os"

	"texApi/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	conn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		config.ENV.DB_USER, config.ENV.DB_PASSWORD,
		config.ENV.DB_HOST, config.ENV.DB_PORT, config.ENV.DB_NAME,
	)

	pool, err := pgxpool.New(context.Background(), conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	DB = pool
}
