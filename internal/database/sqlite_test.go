package database

import (
	"context"
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
		return // Return early to avoid nil pointer dereference
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