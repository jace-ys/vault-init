package storage

import (
	"fmt"

	"github.com/jace-ys/vault-init/pkg/vault"
)

var backends = map[string]BackendInitFunc{
	"postgres": initPostgresStorageBackend,
}

type BackendConfig struct {
	Name     string
	Postgres *PostgresStorage
}

func NewBackendConfig() *BackendConfig {
	return &BackendConfig{
		Postgres: new(PostgresStorage),
	}
}

type BackendInitFunc func(cfg *BackendConfig) (vault.StorageBackend, error)

func UseBackend(cfg *BackendConfig) (vault.StorageBackend, error) {
	initFunc, ok := backends[cfg.Name]
	if !ok {
		return nil, fmt.Errorf("no storage backend named %s", cfg.Name)
	}

	backend, err := initFunc(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize %s storage backend: %s", cfg.Name, err)
	}

	return backend, nil
}

func initPostgresStorageBackend(cfg *BackendConfig) (vault.StorageBackend, error) {
	return NewPostgresStorage(cfg.Postgres.ConnectionURL)
}
