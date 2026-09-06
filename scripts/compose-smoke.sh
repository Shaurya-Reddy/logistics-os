#!/bin/sh
set -eu
mkdir -p reports
docker compose up --build -d --wait --wait-timeout 120
curl --fail --silent http://127.0.0.1:8080/health
curl --fail --silent http://127.0.0.1:8080/ready
curl --fail --silent http://127.0.0.1:8080/ > reports/shell.html
# The running application has no child worker/runtime process.
docker compose exec -T app sh -ec 'test "$(cat /proc/1/comm)" = logistics-os; test -z "$(cat /proc/1/task/1/children)"'
# Runtime role cannot mutate migration history.
if docker compose exec -T postgres sh -ec 'PGPASSWORD="$APP_DB_PASSWORD" psql -h 127.0.0.1 -U logistics_app -d logistics -v ON_ERROR_STOP=1 -c "DELETE FROM schema_migrations"'; then
  echo 'Runtime role unexpectedly has write access' >&2; exit 1
fi
# A bad migration connection must fail, never start serving.
if docker compose run --rm --no-deps -e MIGRATION_DATABASE_URL=postgres://invalid:invalid@postgres/logistics app; then
  echo 'Migration failure did not stop startup' >&2; exit 1
fi
docker compose stop postgres
test "$(curl -s -o reports/not-ready.json -w '%{http_code}' http://127.0.0.1:8080/ready)" = 503
curl --fail --silent http://127.0.0.1:8080/health
docker compose start postgres
attempt=0
until curl --fail --silent http://127.0.0.1:8080/ready; do
  attempt=$((attempt + 1)); test "$attempt" -lt 30; sleep 1
done
