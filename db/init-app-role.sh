#!/bin/sh
set -eu
# psql variable quoting prevents passwords from being interpreted as SQL.
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" --set=app_password="$APP_DB_PASSWORD" <<'SQL'
CREATE ROLE logistics_app LOGIN PASSWORD :'app_password';
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO logistics_app;
ALTER DEFAULT PRIVILEGES FOR ROLE logistics_owner IN SCHEMA public GRANT SELECT ON TABLES TO logistics_app;
SQL
