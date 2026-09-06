-- Only migration bookkeeping exists in S0; no business tables.
CREATE TABLE schema_migrations (
    version text PRIMARY KEY,
    checksum text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT now()
);
