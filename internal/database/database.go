package database

import (
	"context"
	"database/sql"

	"tgbase/config"
)

// FromConfig opens the database specified in cfg.
func FromConfig(ctx context.Context, cfg *config.Config) (Database, error) {
	if cfg.Database.Type == "postgres" {
		return NewPostgresDB(ctx, cfg.Database.Postgres.DSN)
	}
	return NewSQLiteDB(ctx, cfg.Database.SQLite.Path)
}

type Database interface {
	Close() error

	// Core database operations
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row

	// CRUD operations
	Insert(ctx context.Context, table string, data map[string]any) (sql.Result, error)
	Update(ctx context.Context, table string, data map[string]any, where string, args ...any) (sql.Result, error)
	Delete(ctx context.Context, table string, where string, args ...any) (sql.Result, error)
	Select(ctx context.Context, table string, columns []string, where string, args ...any) (*sql.Rows, error)
	SelectRow(ctx context.Context, table string, columns []string, where string, args ...any) *sql.Row
}

// SoftDeleteDatabase extends Database for tables that have a deleted_at column.
type SoftDeleteDatabase interface {
	Database
	// SoftDelete sets deleted_at = now() for matched rows.
	SoftDelete(ctx context.Context, table string, where string, args ...any) (sql.Result, error)
	// Restore clears deleted_at for matched rows.
	Restore(ctx context.Context, table string, where string, args ...any) (sql.Result, error)
	// HardDelete permanently removes rows.
	HardDelete(ctx context.Context, table string, where string, args ...any) (sql.Result, error)
	// SelectDeleted returns only soft-deleted rows.
	SelectDeleted(ctx context.Context, table string, columns []string, where string, args ...any) (*sql.Rows, error)
}
