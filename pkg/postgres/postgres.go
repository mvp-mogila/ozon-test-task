package postgres

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func NewPgxDatabase(ctx context.Context, cfg *PostgresConfig) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		net.JoinHostPort(cfg.Host, cfg.Port),
		cfg.Database,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("error connecting to Postgres")
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("ping Postgres error")
	}
	return pool
}
