package database

import (
	"context"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestNewSQLiteDB(t *testing.T) {
	ctx := context.Background()
	path := "test.db"

	db, err := NewSQLiteDB(ctx, path)
	if err != nil {
		t.Errorf("Failed to create SQLiteDB: %v", err)
	}

	// Add assertions or additional test cases here
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
