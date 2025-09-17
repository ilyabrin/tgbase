package database

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestNewSQLiteDB(t *testing.T) {
	ctx := context.Background()
	path := "test.db"

	db, err := NewSQLiteDB(ctx, path)
	if err != nil {
		t.Errorf("Failed to create SQLiteDB: %v", err)
		return
	}

	if err := db.Close(); err != nil {
		t.Errorf("Failed to close SQLiteDB: %v", err)
	}

	// Clean up the test database file
	if err := os.Remove(path); err != nil {
		t.Errorf("Failed to remove test database file: %v", err)
	}

	// Check if the file was removed
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("Test database file still exists: %v", err)
	}
}

func TestNewSQLiteDB_InvalidPath(t *testing.T) {
	ctx := context.Background()
	// Use an invalid path (directory that doesn't exist)
	path := "/nonexistent/directory/test.db"

	_, err := NewSQLiteDB(ctx, path)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestSQLiteDB_Operations(t *testing.T) {
	ctx := context.Background()
	path := "test_operations.db"

	// Clean up function
	cleanup := func() {
		os.Remove(path)
	}
	defer cleanup()

	db, err := NewSQLiteDB(ctx, path)
	if err != nil {
		t.Fatalf("Failed to create SQLiteDB: %v", err)
	}
	defer db.Close()

	// Create test table
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			age INTEGER,
			email TEXT
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	t.Run("Insert", func(t *testing.T) {
		data := map[string]any{
			"name":  "John Doe",
			"age":   30,
			"email": "john@example.com",
		}

		result, err := db.Insert(ctx, "users", data)
		if err != nil {
			t.Errorf("Insert failed: %v", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			t.Errorf("Failed to get last insert ID: %v", err)
		}

		if lastID <= 0 {
			t.Errorf("Expected positive last insert ID, got %d", lastID)
		}
	})

	t.Run("Query", func(t *testing.T) {
		rows, err := db.Query(ctx, "SELECT id, name, age FROM users WHERE name = ?", "John Doe")
		if err != nil {
			t.Errorf("Query failed: %v", err)
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			var id, age int
			var name string
			if err := rows.Scan(&id, &name, &age); err != nil {
				t.Errorf("Scan failed: %v", err)
			}
			if name != "John Doe" {
				t.Errorf("Expected name 'John Doe', got '%s'", name)
			}
			if age != 30 {
				t.Errorf("Expected age 30, got %d", age)
			}
			count++
		}

		if count != 1 {
			t.Errorf("Expected 1 row, got %d", count)
		}
	})

	t.Run("QueryRow", func(t *testing.T) {
		row := db.QueryRow(ctx, "SELECT name, age FROM users WHERE name = ?", "John Doe")

		var name string
		var age int
		if err := row.Scan(&name, &age); err != nil {
			t.Errorf("QueryRow failed: %v", err)
		}

		if name != "John Doe" {
			t.Errorf("Expected name 'John Doe', got '%s'", name)
		}
		if age != 30 {
			t.Errorf("Expected age 30, got %d", age)
		}
	})

	t.Run("Update", func(t *testing.T) {
		data := map[string]any{
			"age": 31,
		}

		result, err := db.Update(ctx, "users", data, "name = ?", "John Doe")
		if err != nil {
			t.Errorf("Update failed: %v", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Errorf("Failed to get rows affected: %v", err)
		}

		if rowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", rowsAffected)
		}

		// Verify update
		row := db.QueryRow(ctx, "SELECT age FROM users WHERE name = ?", "John Doe")
		var age int
		if err := row.Scan(&age); err != nil {
			t.Errorf("Failed to verify update: %v", err)
		}
		if age != 31 {
			t.Errorf("Expected age 31 after update, got %d", age)
		}
	})

	t.Run("Select", func(t *testing.T) {
		rows, err := db.Select(ctx, "users", []string{"name", "age"}, "age > ?", 25)
		if err != nil {
			t.Errorf("Select failed: %v", err)
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			var name string
			var age int
			if err := rows.Scan(&name, &age); err != nil {
				t.Errorf("Scan failed: %v", err)
			}
			count++
		}

		if count != 1 {
			t.Errorf("Expected 1 row, got %d", count)
		}
	})

	t.Run("SelectRow", func(t *testing.T) {
		row := db.SelectRow(ctx, "users", []string{"name", "email"}, "age = ?", 31)

		var name, email string
		if err := row.Scan(&name, &email); err != nil {
			t.Errorf("SelectRow failed: %v", err)
		}

		if name != "John Doe" {
			t.Errorf("Expected name 'John Doe', got '%s'", name)
		}
		if email != "john@example.com" {
			t.Errorf("Expected email 'john@example.com', got '%s'", email)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		result, err := db.Delete(ctx, "users", "name = ?", "John Doe")
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Errorf("Failed to get rows affected: %v", err)
		}

		if rowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", rowsAffected)
		}

		// Verify deletion
		row := db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE name = ?", "John Doe")
		var count int
		if err := row.Scan(&count); err != nil {
			t.Errorf("Failed to verify deletion: %v", err)
		}
		if count != 0 {
			t.Errorf("Expected 0 rows after deletion, got %d", count)
		}
	})
}

func TestSQLiteDB_ErrorCases(t *testing.T) {
	ctx := context.Background()
	path := "test_errors.db"

	cleanup := func() {
		os.Remove(path)
	}
	defer cleanup()

	db, err := NewSQLiteDB(ctx, path)
	if err != nil {
		t.Fatalf("Failed to create SQLiteDB: %v", err)
	}
	defer db.Close()

	t.Run("Query NonExistent Table", func(t *testing.T) {
		_, err := db.Query(ctx, "SELECT * FROM nonexistent_table")
		if err == nil {
			t.Error("Expected error for nonexistent table, got nil")
		}
	})

	t.Run("Insert Invalid Data", func(t *testing.T) {
		// Try to insert into a table that doesn't exist
		data := map[string]any{
			"name": "test",
		}
		_, err := db.Insert(ctx, "nonexistent_table", data)
		if err == nil {
			t.Error("Expected error for nonexistent table, got nil")
		}
	})

	t.Run("QueryRow No Rows", func(t *testing.T) {
		// Create a table first
		_, err := db.Exec(ctx, "CREATE TABLE IF NOT EXISTS empty_table (id INTEGER)")
		if err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}

		row := db.QueryRow(ctx, "SELECT id FROM empty_table WHERE id = ?", 999)
		var id int
		err = row.Scan(&id)
		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows, got %v", err)
		}
	})
}