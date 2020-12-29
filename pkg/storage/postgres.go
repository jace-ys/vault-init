package storage

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/lib/pq"

	"github.com/jace-ys/vault-init/pkg/vault"
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

func (s *PostgresStorage) WriteInitData(ctx context.Context, data *vault.InitData) error {
	return s.Transact(ctx, func(tx *sql.Tx) error {
		query := `
		INSERT INTO vault_init_data (encryption_type, encryption_version, root_token, unseal_keys)
		VALUES ($1, $2, $3, $4)`
		_, err := tx.ExecContext(ctx, query, data.EncryptionType, data.EncryptionVersion, data.RootToken, pq.Array(data.UnsealKeys))
		return err
	})
}

func (s *PostgresStorage) ReadInitData(ctx context.Context) (*vault.InitData, error) {
	var data vault.InitData
	err := s.Transact(ctx, func(tx *sql.Tx) error {
		query := `
		SELECT encryption_type, encryption_version, root_token, unseal_keys
		FROM vault_init_data
		ORDER BY created_at
		LIMIT 1
		`
		row := tx.QueryRowContext(ctx, query)
		return row.Scan(&data.EncryptionType, &data.EncryptionVersion, &data.RootToken, pq.Array(&data.UnsealKeys))
	})

	return &data, err
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
