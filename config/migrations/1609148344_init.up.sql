CREATE TABLE vault_init_data (
  encryption_type TEXT,
  encryption_version TEXT,
  root_token TEXT,
  unseal_keys TEXT[],
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
