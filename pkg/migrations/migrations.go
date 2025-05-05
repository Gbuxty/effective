package migrations

import (
	"Effective/pkg/logger"
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "./migrations"
	dialect       = "postgres"
)

func Run(ctx context.Context, db *pgxpool.Pool, logger *logger.Logger) error {
	logger.Info("Running migrations...")

	conn, err := db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	goose.SetBaseFS(os.DirFS(migrationsDir))

	dbSQL := stdlib.OpenDBFromPool(db)
	defer dbSQL.Close()
	//хард код как-то убрать
	if err := goose.Up(dbSQL, "."); err != nil {
		return err
	}

	logger.Info("Migrations completed successfully")
	return nil
}
