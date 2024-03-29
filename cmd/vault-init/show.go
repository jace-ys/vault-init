package main

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jace-ys/vault-init/pkg/encryption"
	"github.com/jace-ys/vault-init/pkg/storage"
)

type ShowCommand struct {
	Timeout time.Duration

	Encryption *encryption.BackendConfig
	Storage    *storage.BackendConfig
}

func attachShowCommand(cmd *kingpin.CmdClause) *ShowCommand {
	show := &ShowCommand{
		Encryption: encryption.NewBackendConfig(),
		Storage:    storage.NewBackendConfig(),
	}

	cmd.Flag("timeout", "The duration to wait before timing out the process.").
		Envar("TIMEOUT").
		Default("5s").
		DurationVar(&show.Timeout)

	cmd.Flag("encryption", "The encryption backend to use. Must be one of [local].").
		Envar("ENCRYPTION").
		Default("local").
		EnumVar(&show.Encryption.Name, "local")

	cmd.Flag("encryption-local-secret-key", "The 32-byte secret key to use for encrypting root tokens and unseal keys.").
		Envar("ENCRYPTION_LOCAL_SECRET_KEY").
		PlaceHolder("SECRET-KEY").
		StringVar(&show.Encryption.Local.SecretKey)

	cmd.Flag("storage", "The storage backend to use. Must be one of [postgres].").
		Envar("STORAGE").
		Default("postgres").
		EnumVar(&show.Storage.Name, "postgres")

	cmd.Flag("storage-postgres-connection-url", "The URL to use for connecting to the Postgres server.").
		Envar("STORAGE_POSTGRES_CONNECTION_URL").
		PlaceHolder("CONNECTION-URL").
		URLVar(&show.Storage.Postgres.ConnectionURL)

	return show
}

func (c *ShowCommand) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	encryptionBackend, err := encryption.UseBackend(c.Encryption)
	if err != nil {
		return err
	}

	storageBackend, err := storage.UseBackend(c.Storage)
	if err != nil {
		return err
	}

	data, err := storageBackend.ReadInitData(ctx)
	if err != nil {
		return err
	}

	rootToken, err := encryptionBackend.Decrypt(ctx, data.RootToken)
	if err != nil {
		return err
	}

	unsealKeys := make([]string, len(data.UnsealKeys))
	for i, encryptedUnsealKey := range data.UnsealKeys {
		decrypted, err := encryptionBackend.Decrypt(ctx, encryptedUnsealKey)
		if err != nil {
			return err
		}
		unsealKeys[i] = decrypted
	}

	fmt.Printf("Initial Root Token: %s\n", rootToken)
	for i, unsealKey := range unsealKeys {
		fmt.Printf("Unseal Key %d: %s\n", i+1, unsealKey)
	}

	return nil
}
