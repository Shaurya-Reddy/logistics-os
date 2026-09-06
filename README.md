# Logistics OS

A lightweight, open-source, self-hosted operating system for logistics operations.

## Current milestone

S0 technical skeleton: Go server, Svelte/TypeScript/Vite shell, PostgreSQL, SQL migrations, health/readiness endpoints, Docker Compose, OpenAPI starter contract, and initial CI budget checks.

Business modules are intentionally not implemented yet.

## Run with Docker

```bash
docker compose up --build
```

Then open `http://localhost:8080`.

System endpoints:

- `GET /health` — process liveness; does not require a successful database ping.
- `GET /ready` — database readiness.

The container runs pending SQL migrations before starting the HTTP server.

## Development

Backend:

```bash
export DATABASE_URL='postgres://logistics:logistics_dev@localhost:5432/logistics_os?sslmode=disable'
go run ./cmd/logistics-os migrate
go run ./cmd/logistics-os serve
```

Frontend:

```bash
cd web
npm install
npm run dev
```

The Vite development server proxies `/health` and `/ready` to the Go server on port 8080.

## Project rules

Read `PRODUCT.md`, `ARCHITECTURE.md`, `AGENTS.md`, `PERFORMANCE_BUDGET.md`, and `docs/ERP_REFERENCE.md` before substantial changes.
