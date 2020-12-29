package storage

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/lib/pq"
)

type PostgresStorage struct {
	ConnectionURL *url.URL

	db *sql.DB
}

func NewPostgresStorage(connectionURL *url.URL) (*PostgresStorage, error) {
	if connectionURL == nil {
		return nil, fmt.Errorf("no connection URL provided")
	}

	db, err := sql.Open("postgres", connectionURL.String())
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) WriteInitData(ctx context.Context, token string, keys []string) error {
	return s.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "INSERT INTO vault_init_data (root_token, unseal_keys) VALUES ($1, $2)", token, pq.Array(keys))
		return err
	})
}

func (s *PostgresStorage) ReadInitData(ctx context.Context) (string, []string, error) {
	var token string
	var keys []string
	err := s.Transact(ctx, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, "SELECT root_token, unseal_keys FROM vault_init_data ORDER BY created_at LIMIT 1")
		return row.Scan(&token, pq.Array(&keys))
	})

	return token, keys, err
}

func (s *PostgresStorage) Transact(ctx context.Context, txFunc func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres transaction failed: %w", err)
	}

	err = txFunc(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("postgres transaction failed: %w", err)
	}

	return tx.Commit()
}
