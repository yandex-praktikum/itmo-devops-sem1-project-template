package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"project_sem/internal/model"
	"project_sem/internal/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

//nolint:gochecknoglobals // Настройка, переиспользуемая в рамках пакета
var psql = getQueryBuilder()

type MarketingRepository struct {
	postgresPool *storage.PostgresPool
}

func NewMarketingRepository(postgresPool *storage.PostgresPool) *MarketingRepository {
	return &MarketingRepository{postgresPool: postgresPool}
}

func (r *MarketingRepository) UploadProducts(ctx context.Context, products []model.Product) error {
	tx, err := r.postgresPool.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire connection for transaction %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Info(fmt.Sprintf("failed to rollback transaction: %s", err.Error()))
		}
	}()

	insertProductQuery := `insert into "products" (id, name, category, price, create_date) values ($1, $2, $3, $4, $5)`

	productStatement, err := tx.Prepare(ctx, "insertproduct", insertProductQuery)
	if err != nil {
		return fmt.Errorf("unable to prepare query %w", err)
	}

	batch := &pgx.Batch{}

	for i := 0; i < len(products); i++ {
		batch.Queue(productStatement.Name, products[i].ID, products[i].Name, products[i].Category, products[i].Price, products[i].CreateDate)
	}

	result := tx.SendBatch(ctx, batch)

	if err := result.Close(); err != nil {
		return fmt.Errorf("error executing batch: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("unable to commit %w", err)
	}

	return nil
}

func (r *MarketingRepository) LoadProducts(ctx context.Context) ([]model.Product, error) {
	queryBuilder := psql.Select(
		"id",
		"name",
		"category",
		"price",
		"create_date",
	).From("products")

	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.postgresPool.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get datasources list: %w", err)
	}

	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Product])
	if err != nil {
		return nil, fmt.Errorf("collect product rows: %w", err)
	}

	return products, nil
}

func getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
