# Storage Backends

Storage backends are used to store encrypted root tokens and unseal keys. The same backend should always be used for each initialization of Vault, as each backend can only read data written by itself. These backends are to be specified via the `vault-init` CLI, using the `--storage` flag followed by its name (eg. `--storage=postgres`).

The implementation of these storage backends can be found under [`pkg/storage`](../pkg/storage).

## PostgreSQL (`postgres`)

The `postgres` storage backend reads and writes encrypted data using a PostgreSQL table with the name `vault_init_data`. To use this backend, you will need to create the following table in your PostgreSQL database:

```sql
CREATE TABLE vault_init_data (
  encryption_type TEXT,
  encryption_version TEXT,
  root_token TEXT,
  unseal_keys TEXT[],
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

#### Configuration

- `--postgres-storage-connection-URL`: The URL to use for connecting to the Postgres server.
