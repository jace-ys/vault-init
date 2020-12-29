CREATE TABLE vault_init_data (
  root_token TEXT,
  unseal_keys TEXT[],
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
