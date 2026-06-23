package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMock(t *testing.T) (*PostgresDB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return &PostgresDB{db: db}, mock
}

func TestPostgresDB_ErrorCases(t *testing.T) {
	ctx := context.Background()

	t.Run("Insert Error", func(t *testing.T) {
		pdb, mock := newMock(t)
		mock.ExpectExec("INSERT INTO users").
			WithArgs("john", 25).
			WillReturnError(fmt.Errorf("insert error"))

		data := map[string]any{"name": "john", "age": 25}
		_, err := pdb.Insert(ctx, "users", data)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Update Error", func(t *testing.T) {
		pdb, mock := newMock(t)
		mock.ExpectExec("UPDATE users").
			WithArgs(30, 1).
			WillReturnError(fmt.Errorf("update error"))

		_, err := pdb.Update(ctx, "users", map[string]any{"age": 30}, "id = $2", 1)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Delete Error", func(t *testing.T) {
		pdb, mock := newMock(t)
		mock.ExpectExec("DELETE FROM users").
			WithArgs(1).
			WillReturnError(fmt.Errorf("delete error"))

		_, err := pdb.Delete(ctx, "users", "id = $1", 1)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Select Error", func(t *testing.T) {
		pdb, mock := newMock(t)
		mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(1).
			WillReturnError(fmt.Errorf("select error"))

		_, err := pdb.Select(ctx, "users", []string{"id", "name"}, "id = $1", 1)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Update", func(t *testing.T) {
		pdb, mock := newMock(t)
		mock.ExpectExec("UPDATE users").
			WithArgs(30, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := pdb.Update(ctx, "users", map[string]any{"age": 30}, "id = $2", 1)
		if err != nil {
			t.Errorf("error executing update: %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		pdb, mock := newMock(t)
		mock.ExpectExec("DELETE FROM users").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := pdb.Delete(ctx, "users", "id = $1", 1)
		if err != nil {
			t.Errorf("error executing delete: %v", err)
		}
	})

	t.Run("Select", func(t *testing.T) {
		pdb, mock := newMock(t)
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "john")
		mock.ExpectQuery("SELECT id, name FROM users").
			WithArgs(1).
			WillReturnRows(rows)

		_, err := pdb.Select(ctx, "users", []string{"id", "name"}, "id = $1", 1)
		if err != nil {
			t.Errorf("error executing select: %v", err)
		}
	})

	t.Run("SelectRow", func(t *testing.T) {
		pdb, mock := newMock(t)
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "john")
		mock.ExpectQuery("SELECT id, name FROM users").
			WithArgs(1).
			WillReturnRows(rows)

		row := pdb.SelectRow(ctx, "users", []string{"id", "name"}, "id = $1", 1)
		if row == nil {
			t.Error("expected row, got nil")
		}
	})
}
