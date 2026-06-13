package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(ctx context.Context, dsn string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

// Core helpers
func (p *PostgresDB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

func (p *PostgresDB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

func (p *PostgresDB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

// CRUD methods
func (p *PostgresDB) Insert(ctx context.Context, table string, data map[string]any) (sql.Result, error) {
	var columns []string
	var values []any
	var placeholders []string
	i := 1
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		i++
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","))
	return p.Exec(ctx, query, values...)
}

func (p *PostgresDB) Update(ctx context.Context, table string, data map[string]any, where string, args ...any) (sql.Result, error) {
	var setClauses []string
	var values []any
	i := 1
	for k, v := range data {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", k, i))
		values = append(values, v)
		i++
	}
	values = append(values, args...)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(setClauses, ","),
		where)
	return p.Exec(ctx, query, values...)
}

func (p *PostgresDB) Delete(ctx context.Context, table string, where string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	return p.Exec(ctx, query, args...)
}

func (p *PostgresDB) Select(ctx context.Context, table string, columns []string, where string, args ...any) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "),
		table,
		where)
	return p.Query(ctx, query, args...)
}

func (p *PostgresDB) SelectRow(ctx context.Context, table string, columns []string, where string, args ...any) *sql.Row {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "),
		table,
		where)
	return p.QueryRow(ctx, query, args...)
}

// SoftDeleteDatabase implementation

func (p *PostgresDB) SoftDelete(ctx context.Context, table string, where string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET deleted_at = NOW() WHERE %s AND deleted_at IS NULL", table, where)
	return p.Exec(ctx, query, args...)
}

func (p *PostgresDB) Restore(ctx context.Context, table string, where string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET deleted_at = NULL WHERE %s", table, where)
	return p.Exec(ctx, query, args...)
}

func (p *PostgresDB) HardDelete(ctx context.Context, table string, where string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	return p.Exec(ctx, query, args...)
}

func (p *PostgresDB) SelectDeleted(ctx context.Context, table string, columns []string, where string, args ...any) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s AND deleted_at IS NOT NULL",
		strings.Join(columns, ", "), table, where)
	return p.Query(ctx, query, args...)
}
