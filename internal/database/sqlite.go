package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(ctx context.Context, path string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// Core helpers
func (s *SQLiteDB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *SQLiteDB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *SQLiteDB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

// CRUD methods
func (s *SQLiteDB) Insert(ctx context.Context, table string, data map[string]any) (sql.Result, error) {
	var columns []string
	var values []any
	var placeholders []string
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
		placeholders = append(placeholders, "?")
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","))
	return s.Exec(ctx, query, values...)
}

func (s *SQLiteDB) Update(ctx context.Context, table string, data map[string]any, where string, args ...any) (sql.Result, error) {
	var setClauses []string
	var values []any
	for k, v := range data {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", k))
		values = append(values, v)
	}
	values = append(values, args...)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(setClauses, ","),
		where)
	return s.Exec(ctx, query, values...)
}

func (s *SQLiteDB) Delete(ctx context.Context, table string, where string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	return s.Exec(ctx, query, args...)
}

func (s *SQLiteDB) Select(ctx context.Context, table string, columns []string, where string, args ...any) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "),
		table,
		where)
	return s.Query(ctx, query, args...)
}

func (s *SQLiteDB) SelectRow(ctx context.Context, table string, columns []string, where string, args ...any) *sql.Row {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "),
		table,
		where)
	return s.QueryRow(ctx, query, args...)
}

// SoftDeleteDatabase implementation

func (s *SQLiteDB) SoftDelete(ctx context.Context, table string, where string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET deleted_at = datetime('now') WHERE %s AND deleted_at IS NULL", table, where)
	return s.Exec(ctx, query, args...)
}

func (s *SQLiteDB) Restore(ctx context.Context, table string, where string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET deleted_at = NULL WHERE %s", table, where)
	return s.Exec(ctx, query, args...)
}

func (s *SQLiteDB) HardDelete(ctx context.Context, table string, where string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	return s.Exec(ctx, query, args...)
}

func (s *SQLiteDB) SelectDeleted(ctx context.Context, table string, columns []string, where string, args ...any) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s AND deleted_at IS NOT NULL",
		strings.Join(columns, ", "), table, where)
	return s.Query(ctx, query, args...)
}
