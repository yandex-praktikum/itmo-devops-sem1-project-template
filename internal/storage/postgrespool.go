package storage

import (
	"context"
	"fmt"
	"net"

	"project_sem/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPool struct {
	DB *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*PostgresPool, error) {
	hostPort := net.JoinHostPort(cfg.DBHost, cfg.DBPort)

	dbPool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, hostPort, cfg.DBName))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %w", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database %w", err)
	}

	return &PostgresPool{
		DB: dbPool,
	}, nil
}
