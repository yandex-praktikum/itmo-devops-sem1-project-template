package storage

import (
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Migrations struct {
	db *pgxpool.Pool
}

func NewMigrations(db *pgxpool.Pool) (*Migrations, error) {
	err := goose.SetDialect("pgx")
	if err != nil {
		return nil, fmt.Errorf("failed to set dialect: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	return &Migrations{
		db: db,
	}, nil
}

func (m *Migrations) Up() error {
	// Получаем стандартный *sql.DB из pgxpool.Pool
	db := stdlib.OpenDBFromPool(m.db)
	defer db.Close()

	err := goose.Up(db, "migrations")
	if err != nil {
		return fmt.Errorf("failed to up migrations: %w", err)
	}

	return nil
}
