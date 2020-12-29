package encryption

import (
	"fmt"

	"github.com/jace-ys/vault-init/pkg/vault"
)

var backends = map[string]BackendInitFunc{
	"local": initLocalEncryptionBackend,
}

type BackendConfig struct {
	Name  string
	Local *LocalEncryption
}

func NewBackendConfig() *BackendConfig {
	return &BackendConfig{
		Local: new(LocalEncryption),
	}
}

type BackendInitFunc func(cfg *BackendConfig) (vault.EncryptionBackend, error)

func UseBackend(cfg *BackendConfig) (vault.EncryptionBackend, error) {
	initFunc, ok := backends[cfg.Name]
	if !ok {
		return nil, fmt.Errorf("no encryption backend named %s", cfg.Name)
	}

	backend, err := initFunc(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize %s encryption backend: %s", cfg.Name, err)
	}

	return backend, nil
}

func initLocalEncryptionBackend(cfg *BackendConfig) (vault.EncryptionBackend, error) {
	return NewLocalEncryption(cfg.Local.SecretKey)
}
