CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(19) PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
