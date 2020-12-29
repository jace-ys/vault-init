package vault

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type VaultClient interface {
	Health(ctx context.Context) (*http.Response, error)
	Init(ctx context.Context, opts *InitRequest) (*InitResponse, error)
	Unseal(ctx context.Context, opts *UnsealRequest) (*UnsealResponse, error)
}

type EncryptionBackend interface {
	Encrypt(ctx context.Context, plaintext string) (string, error)
	Decrypt(ctx context.Context, data string) (string, error)
}

type StorageBackend interface {
	WriteInitData(ctx context.Context, token string, keys []string) error
	ReadInitData(ctx context.Context) (string, []string, error)
}

type VaultInit struct {
	client     VaultClient
	encryption EncryptionBackend
	storage    StorageBackend

	once          bool
	checkInterval time.Duration
}

func NewVaultInit(client VaultClient, encryption EncryptionBackend, storage StorageBackend, once bool, checkInterval time.Duration) *VaultInit {
	return &VaultInit{
		client:        client,
		encryption:    encryption,
		storage:       storage,
		once:          once,
		checkInterval: checkInterval,
	}
}

func (vi *VaultInit) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			return nil
		default:
			response, err := vi.client.Health(ctx)
			if err != nil {
				log.Printf("Failed to retrieve Vault status: %s\n", err)
			} else {
				switch response.StatusCode {
				case 200:
					log.Println("Vault is initialized, unsealed, and active")
				case 429:
					log.Println("Vault is unsealed and in standby")
				case 501:
					vi.handle501(ctx)
				case 503:
					vi.handle503(ctx)
				default:
					log.Printf("Vault status is unknown (response code = %d)", response.StatusCode)
				}
			}
		}

		if vi.once {
			return nil
		}

		log.Println("Sleeping until next health check...")

		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			return nil
		case <-time.After(vi.checkInterval):
			continue
		}
	}
}

func (vi *VaultInit) handle501(ctx context.Context) {
	log.Println("Vault is not initialized; attempting to initialize and unseal...")

	err := vi.initialize(ctx)
	if err != nil {
		log.Printf("Failed to initialize Vault: %s\n", err)
		return
	}

	log.Println("Vault was successfully initialized")

	err = vi.unseal(ctx)
	if err != nil {
		log.Printf("Failed to unseal Vault: %s\n", err)
		return
	}

	log.Println("Vault was successfully unsealed")
}

func (vi *VaultInit) handle503(ctx context.Context) {
	log.Println("Vault is sealed; attempting to unseal...")

	err := vi.unseal(ctx)
	if err != nil {
		log.Printf("Failed to unseal Vault: %s\n", err)
		return
	}

	log.Println("Vault was successfully unsealed")
}

func (vi *VaultInit) initialize(ctx context.Context) error {
	opts := &InitRequest{
		SecretShares:    5,
		SecretThreshold: 3,
	}

	initResponse, err := vi.client.Init(ctx, opts)
	if err != nil {
		return err
	}

	encryptedRootToken, err := vi.encryption.Encrypt(ctx, initResponse.RootToken)
	if err != nil {
		return err
	}

	encryptedUnsealKeys := make([]string, len(initResponse.KeysBase64))
	for i, unsealKey := range initResponse.KeysBase64 {
		encrypted, err := vi.encryption.Encrypt(ctx, unsealKey)
		if err != nil {
			return err
		}
		encryptedUnsealKeys[i] = encrypted
	}

	return vi.storage.WriteInitData(ctx, encryptedRootToken, encryptedUnsealKeys)
}

func (vi *VaultInit) unseal(ctx context.Context) error {
	_, encryptedUnsealKeys, err := vi.storage.ReadInitData(ctx)
	if err != nil {
		return err
	}

	if len(encryptedUnsealKeys) == 0 {
		return fmt.Errorf("no unseal keys could be found")
	}

	unsealKeys := make([]string, len(encryptedUnsealKeys))
	for i, encryptedUnsealKey := range encryptedUnsealKeys {
		decrypted, err := vi.encryption.Decrypt(ctx, encryptedUnsealKey)
		if err != nil {
			return err
		}
		unsealKeys[i] = decrypted
	}

	for _, unsealKey := range unsealKeys {
		opts := &UnsealRequest{
			Key: unsealKey,
		}

		unsealResponse, err := vi.client.Unseal(ctx, opts)
		if err != nil {
			return err
		}

		if !unsealResponse.Sealed {
			return nil
		}
	}

	return nil
}
