package database

import (
	"context"
	"testing"

	"tgbase/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSoftDeleteDB(t *testing.T) SoftDeleteDatabase {
	t.Helper()
	ctx := context.Background()
	db, err := NewSQLiteDB(ctx, ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(ctx, `
		CREATE TABLE items (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			deleted_at DATETIME
		)`)
	require.NoError(t, err)

	return db // *SQLiteDB implements SoftDeleteDatabase directly
}

func seedItem(t *testing.T, db SoftDeleteDatabase, name string) int64 {
	t.Helper()
	res, err := db.Insert(context.Background(), "items", map[string]any{"name": name})
	require.NoError(t, err)
	id, err := res.LastInsertId()
	require.NoError(t, err)
	return id
}

func TestSoftDelete(t *testing.T) {
	db := newSoftDeleteDB(t)
	ctx := context.Background()

	id := seedItem(t, db, "Alice")

	res, err := db.SoftDelete(ctx, "items", "id = ?", id)
	require.NoError(t, err)
	n, _ := res.RowsAffected()
	assert.Equal(t, int64(1), n)

	// Row still exists but deleted_at is set.
	var deletedAt *string
	row := db.QueryRow(ctx, "SELECT deleted_at FROM items WHERE id = ?", id)
	require.NoError(t, row.Scan(&deletedAt))
	assert.NotNil(t, deletedAt, "deleted_at should be set after SoftDelete")
}

func TestSoftDelete_Idempotent(t *testing.T) {
	db := newSoftDeleteDB(t)
	ctx := context.Background()

	id := seedItem(t, db, "Bob")
	db.SoftDelete(ctx, "items", "id = ?", id)

	// Second call should match 0 rows (WHERE deleted_at IS NULL).
	res, err := db.SoftDelete(ctx, "items", "id = ?", id)
	require.NoError(t, err)
	n, _ := res.RowsAffected()
	assert.Equal(t, int64(0), n)
}

func TestRestore(t *testing.T) {
	db := newSoftDeleteDB(t)
	ctx := context.Background()

	id := seedItem(t, db, "Carol")
	db.SoftDelete(ctx, "items", "id = ?", id)

	res, err := db.Restore(ctx, "items", "id = ?", id)
	require.NoError(t, err)
	n, _ := res.RowsAffected()
	assert.Equal(t, int64(1), n)

	var deletedAt *string
	row := db.QueryRow(ctx, "SELECT deleted_at FROM items WHERE id = ?", id)
	require.NoError(t, row.Scan(&deletedAt))
	assert.Nil(t, deletedAt, "deleted_at should be NULL after Restore")
}

func TestHardDelete(t *testing.T) {
	db := newSoftDeleteDB(t)
	ctx := context.Background()

	id := seedItem(t, db, "Dave")

	res, err := db.HardDelete(ctx, "items", "id = ?", id)
	require.NoError(t, err)
	n, _ := res.RowsAffected()
	assert.Equal(t, int64(1), n)

	var count int
	row := db.QueryRow(ctx, "SELECT COUNT(*) FROM items WHERE id = ?", id)
	require.NoError(t, row.Scan(&count))
	assert.Equal(t, 0, count, "row should be gone after HardDelete")
}

func TestSelectDeleted(t *testing.T) {
	db := newSoftDeleteDB(t)
	ctx := context.Background()

	id1 := seedItem(t, db, "Eve")
	id2 := seedItem(t, db, "Frank")
	seedItem(t, db, "Grace") // not deleted

	db.SoftDelete(ctx, "items", "id = ?", id1)
	db.SoftDelete(ctx, "items", "id = ?", id2)

	rows, err := db.SelectDeleted(ctx, "items", []string{"id", "name"}, "1=1")
	require.NoError(t, err)
	defer rows.Close()

	found := map[int64]string{}
	for rows.Next() {
		var id int64
		var name string
		require.NoError(t, rows.Scan(&id, &name))
		found[id] = name
	}
	require.NoError(t, rows.Err())

	assert.Len(t, found, 2)
	assert.Equal(t, "Eve", found[id1])
	assert.Equal(t, "Frank", found[id2])
}

func TestSoftDeleteDatabase_Interface(t *testing.T) {
	// Verify both implementations satisfy the interface at compile time.
	var _ SoftDeleteDatabase = &SQLiteDB{}
	var _ SoftDeleteDatabase = &PostgresDB{}
}

// --- FromConfig ---

func TestFromConfig_SQLite(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Type = "sqlite"
	cfg.Database.SQLite.Path = ":memory:"

	db, err := FromConfig(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	_, err = db.Exec(context.Background(), "CREATE TABLE test (id INTEGER PRIMARY KEY)")
	assert.NoError(t, err)
}

func TestFromConfig_DefaultsToSQLite(t *testing.T) {
	// Any type other than "postgres" should open SQLite.
	cfg := &config.Config{}
	cfg.Database.Type = "unknown"
	cfg.Database.SQLite.Path = ":memory:"

	db, err := FromConfig(context.Background(), cfg)
	require.NoError(t, err)
	defer db.Close()
}
