# Architecture foundation

Status: required design for incremental implementation, not a description of completed code.
Companion documents: [product](PRODUCT.md), [agent rules](AGENTS.md), [performance budgets](PERFORMANCE_BUDGET.md).

## Fixed stack and runtime

- Backend: Go, using the standard library where practical.
- Frontend: Svelte + TypeScript + Vite, built as static assets and served by Go, preferably through go:embed.
- Persistence: PostgreSQL; pgx for connections and explicit SQL; sqlc generates typed access from reviewed SQL. No ORM or automatic schema synchronization.
- Interface: versioned REST/JSON with an OpenAPI contract.
- Structure: modular monolith, one deployable application binary/container plus PostgreSQL. Background workers run inside that application process.
- Node is a build/development tool only, never a production server or runtime image requirement.

Do not add Redis, Kafka, RabbitMQ, Elasticsearch, Kubernetes, microservices, or another mandatory service without measured evidence, an architecture decision, and explicit maintainer approval. First try indexes, SQL, bounded in-process work, and PostgreSQL-backed jobs. No distributed system by default.

## Module boundaries and proposed layout

Use cmd/logistics-os for the executable, internal/<module> for business code, internal/platform for shared infrastructure, db/migrations for ordered SQL, module-local SQL queries and generated sqlc access, api/openapi.yaml for the contract, and web/ for the Svelte client. This is a target layout; do not create empty module scaffolding.

Module ownership:
- Identity: users, credentials, sessions, memberships, permission grants.
- Organization: organizations, sites, and organization settings.
- Customers: customer master data and customer-owned item/SKU definitions.
- Warehouse: warehouses, locations, inbound receipts and putaway orchestration.
- Inventory: immutable stock movements, current balance projections, and stock invariants.
- Audit: append-only business audit records and scoped read access.
- Platform: database wiring, HTTP concerns, configuration, job execution, file adapters, logging; no business policy.

Future modules include transport, finance/ERP, quality, CRM expansion, and workforce. Add only when requested. Localization packages provide jurisdiction-specific policy and presentation without importing that policy into the neutral core.

Each module owns writes to its tables. Invoke another module's explicit Go service interface rather than writing its tables or importing its internal generated queries. Keep dependencies acyclic; application wiring coordinates modules. Warehouse calls Inventory to post movements; Inventory does not depend on receipt handlers. A receipt post, balance updates, audit, and durable job/outbox records share one PostgreSQL transaction through an explicit transaction-scoped service boundary. Cross-module read reports use named, reviewed queries; never hide writes in reporting code. Avoid generic repositories, universal entities, service locators, and speculative event frameworks.

## Data model and integrity

Use stable UUID identifiers generated in the application and explicit foreign keys. Organization-owned tables carry organization_id. Composite keys/foreign keys must prevent references crossing organization boundaries; every request and query also enforces membership and scope. A single-organization installation still uses this model. PostgreSQL row-level security is optional defense in depth, not a replacement for authorization.

A site belongs to an organization; a warehouse belongs to a site. Customers belong to organizations. Locations belong to warehouses; items belong to customers. Include organization/customer ownership in stock keys and constraints. Use unique human-facing codes within a documented scope; codes are not database identities.

Use timestamptz for instants, UTC internally, and IANA time zones for site display. Keep local calendar dates as date values. Quantities and money use exact numeric/decimal representations, never binary floating point. v0.1 quantities use numeric(20,6); reject excess precision rather than silently rounding. Money, when introduced, carries an explicit currency and rounding policy. JSON quantities use decimal strings to preserve precision. Do not assume one currency, address format, or unit system globally.

Use constraints for positive receipt quantities, nonnegative balances, valid relationships, and uniqueness. Post stock changes atomically, locking affected balance rows in deterministic order and handling creation races. In one transaction, enforce available quantity, append movements, update balances, and record audit. Reversals reference original movements and cannot reverse more than the unreversed quantity. Transfers debit and credit together. Current balances are a projection of immutable movements and must be reconcilable from them.

Use optimistic version checks on editable documents; reject stale updates. Disable referenced master data instead of deleting historical meaning. Posted business facts are immutable; corrections are linked compensating records. Define retention explicitly; never silently cascade-delete stock or audit history.

Migrations are ordered, reviewed SQL tracked in source control. Never alter an already released migration. Run migrations as an explicit app subcommand before serving, with a PostgreSQL advisory lock and clear failure. This is a one-shot use of the app image, not a third running service. Use additive/backfill/cleanup changes for safe upgrades. Destructive changes require approval and a tested backup/restore or forward-repair plan.

## API and frontend

Use /api/v1, resource-oriented endpoints, and explicit business actions such as receive and put-away. The server is authoritative for permissions, validation, state transitions, and quantities. Maintain OpenAPI alongside behavior, including schemas, decimal strings, permissions, errors, pagination, and idempotency.

Return consistent problem details with a stable code, safe message, field errors where relevant, and request ID. Use appropriate 400/401/403/404/409/422/429 responses; never leak stack traces or resource existence across scopes. Bound request bodies, filters, timeouts, and exports. Lists default to 50 items and cap at 100, with stable ordering and cursor pagination where appropriate.

Require an idempotency key for receipt posting, putaway, and reversals. Scope keys to organization, actor, and action, store a request hash and result atomically with the operation, reject mismatched reuse, and replay completed results. Retain key records for the v0.1 operational record lifetime; document a safe retention policy before pruning.

Serve assets and API on the same origin; disable broad CORS. Use hashed long-cache assets and a revalidated application HTML shell. Lazy-load feature screens. Svelte renders and collects input; business rules stay on the server. Avoid client secrets, duplicate permission logic as authority, and new state/UI libraries without need.

## Authentication and authorization

v0.1 uses local accounts and PostgreSQL-backed opaque sessions. Hash passwords with a vetted Argon2id implementation and documented, benchmarked parameters; never design cryptography. Store only hashed session tokens, rotate sessions at login/privilege change, enforce idle and absolute expiry, and revoke on logout or account disablement.

Use Secure, HttpOnly, SameSite cookies over production HTTPS, CSRF protection on state-changing requests, login throttling, and generic authentication errors. Provision the first administrator with an explicit bootstrap command and no default password. Recovery is an audited administrator/CLI operation until a secure email flow is requested; SMTP is not mandatory. External OIDC/SSO may be added later as an optional adapter.

Deny by default. Check action, organization, site, and resource ownership on every endpoint and background action. Hiding a button does not authorize a request. Do not accept organization scope from request data without checking membership. Test cross-organization, cross-site, and privilege escalation failures.

## Audit and observability

Append audit records in the same transaction as consequential business changes: actor or system identity, organization/site scope, action, entity identity, timestamp, request/correlation ID, and safe before/after fields or reason. Record receipt posting, putaway, reversal, configuration and access changes. Log authentication failures separately with bounded retention.

Runtime permissions must prevent updating or deleting movement/audit records; use a distinct migration role. Database administrators remain technically capable of changing data, so do not claim cryptographic tamper-proof storage. Never log passwords, session tokens, secrets, or unnecessary personal data.

Emit structured logs to stdout and expose lightweight health and metrics hooks without requiring an observability stack. /health reports process liveness without a database dependency; /ready checks database connectivity and compatible schema with short timeouts. Do not expose sensitive diagnostics publicly.

## Jobs and integration boundaries

Use synchronous transactions for short business operations. Durable deferred work uses a PostgreSQL jobs table with status, due time, attempts, lease expiry, and idempotency key. Workers claim via SELECT ... FOR UPDATE SKIP LOCKED, use bounded concurrency, recover expired leases, retry with capped backoff, and expose terminal failures for manual retry. Assume at-least-once execution, never exactly once.

When external delivery is introduced, persist an outbox entry in the business transaction and make delivery retry-safe. Do not hold database transactions open during network calls. In-process timers are only for disposable maintenance, not durable business work. v0.1 must not build unused integration machinery.

## Storage, deployment, and recovery

Store relational records in PostgreSQL. If attachments are introduced, use a mounted local data directory by default and keep metadata, ownership, checksum, and storage key in PostgreSQL. Enforce authorization, size/type limits, safe generated paths, and reconciliation for incomplete uploads. Never depend on the container's ephemeral filesystem for durable files. An S3-compatible adapter is optional future work, not a required third service.

Use a multi-stage build: Node builds frontend assets; Go builds the server; a minimal non-root runtime contains the binary, assets, and needed CA/time-zone data. Pin supported toolchains, dependencies, and image versions; lock frontend dependencies. Compose starts exactly app and PostgreSQL with persistent volumes. Configuration comes from validated environment/settings; secrets stay outside source control. TLS may terminate at an existing operator proxy or the app; a bundled proxy must not become mandatory.

Use a bounded pgx pool, graceful shutdown, readiness gating, and timeouts. Back up PostgreSQL and any attachment volume consistently, document restore steps, and exercise restore before release. Upgrade only after backup and migration review; an older binary is not assumed compatible with a newer schema. Define recovery steps for each release. Horizontal scaling and high availability require a separate measured design.

## Localization, customization, and optional AI

Keep translations, number/date display, and locale preferences separate from business identifiers and stored values. Jurisdiction-specific tax, customs, statutory reports, validation, and document numbering belong in versioned localization modules with explicit effective dates and fixtures. Do not encode one country's rules in stock or customer core tables.

Prefer typed configuration and explicit extension interfaces. Future custom fields may use validated, namespaced JSONB for supplemental metadata, never essential stock/accounting invariants. Extensions ship as reviewed Go/Svelte modules in the same build; no arbitrary uploaded code, runtime plugin loader, EAV core, or low-code engine.

AI adapters are opt-in and isolated from core paths, with explicit data disclosure, scoped access, time/cost limits, and human confirmation before applying suggestions. Provider outages must not block operations. No model may bypass normal APIs, permissions, or audit.

## Architecture changes

A change to the fixed stack, deployment services, module ownership, auth/trust boundaries, or performance limits requires a concise architecture decision recording the problem, measured evidence, simpler alternatives, operational/security costs, and migration/reversal plan. Obtain explicit maintainer approval before implementing the exception and update all affected foundation documents. Routine work within these boundaries needs no additional approval.
