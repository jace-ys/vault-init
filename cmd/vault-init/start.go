package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jace-ys/vault-init/pkg/encryption"
	"github.com/jace-ys/vault-init/pkg/storage"
	"github.com/jace-ys/vault-init/pkg/vault"
)

type StartCommand struct {
	VaultAddr     *url.URL
	CheckInterval time.Duration
	Once          bool

	Encryption *encryption.BackendConfig
	Storage    *storage.BackendConfig
}

func attachStartCommand(cmd *kingpin.CmdClause) *StartCommand {
	start := &StartCommand{
		Encryption: encryption.NewBackendConfig(),
		Storage:    storage.NewBackendConfig(),
	}

	cmd.Flag("vault-addr", "Address of the Vault server.").
		Envar("VAULT_ADDR").
		Default("https://127.0.0.1:8200").
		URLVar(&start.VaultAddr)

	cmd.Flag("once", "Run the process once and exit instead of running as a daemon.").
		Envar("ONCE").
		Default("false").
		BoolVar(&start.Once)

	cmd.Flag("check-interval", "The time interval between successive health checks of the Vault server.").
		Envar("CHECK_INTERVAL").
		Default("60s").
		DurationVar(&start.CheckInterval)

	cmd.Flag("encryption", "The encryption backend to use. Must be one of [local].").
		Envar("ENCRYPTION").
		Default("local").
		EnumVar(&start.Encryption.Name, "local")

	cmd.Flag("encryption-local-secret-key", "The 32-byte secret key to use for encrypting root tokens and unseal keys.").
		Envar("ENCRYPTION_LOCAL_SECRET_KEY").
		PlaceHolder("SECRET-KEY").
		StringVar(&start.Encryption.Local.SecretKey)

	cmd.Flag("storage", "The storage backend to use. Must be one of [postgres].").
		Envar("STORAGE").
		Default("postgres").
		EnumVar(&start.Storage.Name, "postgres")

	cmd.Flag("storage-postgres-connection-url", "The URL to use for connecting to the Postgres server.").
		Envar("STORAGE_POSTGRES_CONNECTION_URL").
		PlaceHolder("CONNECTION-URL").
		URLVar(&start.Storage.Postgres.ConnectionURL)

	return start
}

func (c *StartCommand) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		select {
		case <-ctx.Done():
			stop()
		}
	}()

	encryptionBackend, err := encryption.UseBackend(c.Encryption)
	if err != nil {
		return err
	}

	storageBackend, err := storage.UseBackend(c.Storage)
	if err != nil {
		return err
	}

	vaultClient := vault.NewClient(c.VaultAddr)
	vaultInit := vault.NewVaultInit(vaultClient, encryptionBackend, storageBackend, c.Once, c.CheckInterval)

	if err := vaultInit.Start(ctx); err != nil {
		return err
	}

	log.Println("Process terminated")
	return nil
}
