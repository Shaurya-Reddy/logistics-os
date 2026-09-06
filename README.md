# Logistics OS

S0 is a technical skeleton: Go serves an embedded Svelte/TypeScript/Vite status page, PostgreSQL stores versioned migration history, and no business module is implemented. The canonical [ERP reference](docs/ERP_REFERENCE.md) preserves the bounded S0–S4 roadmap.

## Start locally

Requires Docker Engine with Compose v2 (Linux amd64 is the budget target).

```sh
cp .env.example .env
# Replace BOTH placeholders in .env with independent `openssl rand -hex 24` values.
docker compose up --build -d --wait
```

Open http://localhost:8080. `GET /health` returns `{"status":"ok"}` independently of PostgreSQL; `GET /ready` returns `{"status":"ready"}` only for a reachable database with the exact expected migration history, otherwise HTTP 503 with `{"status":"unavailable"}`. Each readiness check has a one-second timeout.

Exactly two production services run: app and PostgreSQL. The app startup command invokes the explicit `migrate` subcommand before replacing itself with `serve`; a migration failure prevents serving. Node is build-only. PostgreSQL data persists in the named volume. The HTTP port is loopback-only for local use. Production deployment needs operator-provided HTTPS, secret management and reviewed credentials; no reverse-proxy service is mandatory.

S0 creates separate migration-owner and read-only application roles. The first PostgreSQL initialization script provisions the application role; changing `.env` later does not rotate existing database passwords. The owner URL is removed from the serving process environment after migration, but remains visible to Docker administrators through container configuration. Separate one-shot migration invocation is available when the operator supplies credentials externally:

```sh
docker compose run --rm --entrypoint /app/logistics-os app migrate
```

The migration command takes `MIGRATION_DATABASE_URL`; serving takes `DATABASE_URL` and optional `HTTP_ADDR` (default `:8080`). Migration timeout is 30 seconds; the pool is capped at four connections. Ordered SQL and checksums are embedded in the binary. All pending migrations and history records commit together under a PostgreSQL transaction advisory lock. Released SQL is immutable. An older binary refuses a newer or modified history. HTTP serving never applies DDL.

`docker compose stop` preserves data; `docker compose start` resumes it. Do not remove the volume to solve migration errors. Take a PostgreSQL backup before future upgrades. S0 has no business data, and a tested release backup/restore procedure remains an S4 gate.

## Develop and verify

Pin Go 1.27.1, Node 24.20.0, sqlc 1.31.1 and the lockfile versions. Build the frontend before Go because its output is embedded.

```sh
npm --prefix web ci
npm --prefix web run check
npm --prefix web run build
go mod download
sqlc generate
go vet ./...
go test -race ./...
go build -trimpath -o bin/logistics-os ./cmd/logistics-os
node scripts/links.mjs
node scripts/inventory.mjs
web/node_modules/.bin/redocly lint api/openapi.yaml
node scripts/budgets.mjs
```

To run integration tests, set `TEST_DATABASE_URL` to a **disposable empty database**. The test refuses an existing migration table, runs competing migrations, retries, changed history, failed DDL rollback, and newer-schema rejection. It intentionally leaves this test database with a future-history row; do not use an application database. Without this variable the database test explicitly skips; that is not an integration pass. CI supplies a fresh PostgreSQL instance.

After Compose is ready, run `cd web && npx --no-install playwright install chromium && npx --no-install playwright test`. Browser tests exercise readiness, loading/error/retry and keyboard focus, and compare actual JS/CSS requests against the full Vite entry import graph. Screenshot and network evidence go in `reports/`. `sh scripts/compose-smoke.sh` additionally tests role restrictions and database loss/recovery, and is intended for disposable development/CI deployments because it stops PostgreSQL briefly.

CI checks generated SQL and lockfile drift, Go/type/contract checks, real migrations, fresh Compose, failure behavior, assets, runtime size/contents, service count, inventory drift and vulnerabilities. Evidence is attached to the exact CI commit. Required status checks/branch protection must be configured separately by a repository administrator; this workflow does not enable protection.

The [performance contract](PERFORMANCE_BUDGET.md) remains unchanged. `reports/budgets.json` uses decimal bytes and gzip level 9, counts the complete initial import graph and all chunks, and records missing measurements as pending. There is only one public S0 shell: login and authenticated screens do not exist yet. No S0 result qualifies future screens. Docker image size/contents require `node scripts/budgets.mjs --image`. Controlled-runner startup/RSS enforcement, compressed image transfer size and comparable readiness latency results remain pending until measured on the specified runner; shared CI and laptop runs cannot certify those budgets. Domain load tests are not applicable until domain endpoints exist.

See [dependency decisions](docs/DEPENDENCIES.md) and [inventory](docs/dependencies.json). No project license has been selected, and public release remains gated on that maintainer decision and security/license review.
Fully Integrated Logistics Management System
