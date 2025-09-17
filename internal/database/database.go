package database

import (
	"context"
	"database/sql"
)

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
