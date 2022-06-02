package store

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

const tableAlreadyExists = "table gcache_cache already exists"

type sqliteStore struct {
	db *sql.DB
}

// SQLiteStore creates a SQLite data store.
func SQLiteStore(ctx context.Context, db *sql.DB) (Store, error) {
	err := createTable(ctx, db)
	if err != nil {
		return nil, err
	}
	return &sqliteStore{db}, nil
}

func (s sqliteStore) Get(ctx context.Context, key string) ([]byte, error) {
	var hexString string
	//goland:noinspection SqlNoDataSourceInspection
	row := s.db.QueryRowContext(ctx, "SELECT data FROM gcache_cache WHERE key = ?", key)
	err := row.Scan(&hexString)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	v, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *sqliteStore) Set(ctx context.Context, key string, data []byte) error {
	hexString := hex.EncodeToString(data)
	//goland:noinspection SqlNoDataSourceInspection
	_, err := s.db.ExecContext(ctx, "INSERT OR REPLACE INTO gcache_cache (key, data) VALUES (?, ?)", key, hexString)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteStore) Delete(ctx context.Context, key string) error {
	//goland:noinspection SqlNoDataSourceInspection
	_, err := s.db.ExecContext(ctx, "DELETE FROM gcache_cache WHERE key = ?", key)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteStore) Clear(ctx context.Context) error {
	//goland:noinspection SqlNoDataSourceInspection
	_, err := s.db.ExecContext(ctx, "DELETE FROM gcache_cache")
	if err != nil {
		return err
	}
	return nil
}

func createTable(ctx context.Context, db *sql.DB) error {
	q := `CREATE TABLE gcache_cache ("key" VARCHAR(64) NOT NULL PRIMARY KEY, "data" TEXT);`

	statement, err := db.PrepareContext(ctx, q)
	if err != nil {
		if err.Error() == tableAlreadyExists {
			return nil
		}
		return err
	}

	_, err = statement.ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
