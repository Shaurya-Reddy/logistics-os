-- name: SchemaHistory :many
SELECT version, checksum FROM schema_migrations ORDER BY version;
