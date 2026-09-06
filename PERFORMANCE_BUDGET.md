# Performance and complexity budgets

Status: acceptance requirements for the future implementation. No benchmark baseline or CI enforcement currently exists merely because this document exists.
Read with [PRODUCT.md](PRODUCT.md), [ARCHITECTURE.md](ARCHITECTURE.md), and [AGENTS.md](AGENTS.md).

## Measurement contract

Use decimal KB/MB (1 KB = 1,000 bytes; 1 MB = 1,000,000 bytes). Measure production builds with pinned toolchains and lockfiles, not development servers. Store raw values, commit SHA, versions, dataset seed, runner CPU/memory/architecture, and exact commands with every report.

Reference deployment: Linux amd64, 2 vCPU and 4 GiB RAM, local SSD, app and PostgreSQL on the same host. Limit the application to 1 vCPU/512 MiB and PostgreSQL to 1 vCPU/2 GiB. Disable development instrumentation; keep production auth, audit, logging, and constraints enabled. Use a separate load-generator runner on the same local network. Record deviations; laptop timings or shared-runner noise cannot establish a comparable baseline.

Budgets are release acceptance limits. Deterministic size/service/dependency checks are blocking in ordinary CI once the skeleton exists. Timing and RAM checks are blocking on a controlled runner; until that runner exists, attach reproducible manual results and mark automated enforcement pending. Missing results are not a pass.

## Frontend JavaScript and CSS

- Initial JavaScript: <= 200 KB gzip.
- Initial CSS: <= 50 KB gzip.
- Total built JavaScript across all application chunks: <= 600 KB gzip.
- Total built CSS across all chunks: <= 100 KB gzip.

Use gzip level 9 with reproducible metadata and sum each unique asset once. The initial budget applies separately to login and the authenticated landing screen, including shared/vendor chunks and all automatically loaded imports or preloads needed to make each screen interactive. Report the larger value. Include any CDN-loaded code; required remote runtime assets are prohibited. Source maps are excluded from transfer totals but must not be public by default.

Measure CSS delivered as assets or inline styles, and count eagerly fetched feature chunks. Lazy loading must defer actual work, not hide imports from the calculation. Future screens must fit the same 200 KB JS/50 KB CSS cold-entry budget for their initial interactive view. Track fonts/images separately and avoid adding large assets without review.

CI must parse the Vite manifest and confirm the expected network requests in a production browser smoke test. Never calculate only the entry file while omitting its dependency graph.

## Application container

Target and release ceiling: <= 100 MB uncompressed cumulative runtime image size for linux/amd64, as reported by docker image inspect .Size. Also record compressed registry transfer size separately; it is not the acceptance metric.

Include binary, embedded/static assets, certificates, time-zone data, and runtime base layers. Exclude build stages and PostgreSQL. No Node runtime, frontend source tree, compiler, or package-manager cache in the final image. CI builds the production image and checks its size and contents. Other release architectures require their own recorded results.

## Idle memory

Target and release ceiling: <= 150 MB application process RSS, excluding PostgreSQL.

On the reference deployment, after readiness and a five-minute idle settling period with only health probes, sample once per second for 60 seconds and use the maximum RSS. Keep normal configured pools and background workers active. Record app cgroup memory and PostgreSQL memory separately; a process RSS pass is not proof of whole-host capacity. CI must detect unexpected child processes; moving work to another process cannot bypass this budget.

## Startup

Target: < 2.0 seconds from container process start to the first successful /ready response.

Measure against an already running, reachable PostgreSQL with migrations already applied. Exclude image pulling, database startup, and migrations from this metric; report end-to-end Compose startup and migration duration separately. Probe every 50 ms; run 20 fresh app starts and require the nearest-rank p95 (19th sorted sample) to be < 2.0 seconds. Record cache conditions. A server answering /health before schema/database readiness does not pass.

## Common API latency

On the reference deployment:
- Scoped detail reads and indexed paginated lists: p95 <= 200 ms, p99 <= 500 ms.
- Transactional writes, including receipt posting and putaway: p95 <= 300 ms, p99 <= 750 ms.
- Unexpected errors: zero during the measured valid-request workload.

These are load-generator-observed HTTP durations on the local network, including auth, validation, SQL, stock invariants, audit, and response encoding. Browser rendering and wide-area latency are excluded. Report each endpoint separately; averages across endpoints cannot conceal a slow operation.

Fixture: two organizations, each with two sites, two warehouses per site, 100 customers, 10,000 items, 100,000 stock balance rows, and 1,000,000 stock movement rows. Seed realistic indexed distributions and history for receipts; use deterministic synthetic data, not production exports.

Use 20 concurrent authenticated users, a two-minute warmup, then ten measured minutes with an 80% read/20% write workload. Lists contain 50 rows; write documents contain up to 20 lines. Include ordinary and contended stock keys, distinct valid operations, and sufficient inventory. Obtain at least 1,000 samples per representative endpoint; extend the run if needed. Report achieved throughput, p50/p95/p99, failures, resource use, and fixture version. Run retry/conflict correctness scenarios separately and classify intentional 409 responses rather than hiding unexpected errors.

At skeleton stage, measure a database-backed readiness probe and mark domain benchmarks not applicable. Introduce fixture subsets as endpoints arrive; full v0.1 qualification requires the fixture and endpoint coverage above. Bulk exports and future reports need separate explicit limits and bounded jobs before implementation; they are not excuses for unbounded common APIs.

## Mandatory services

Exactly two mandatory long-running production services: application and PostgreSQL. The worker runs inside the app; migrations use a one-shot invocation of the same image.

A default installation must need no Redis, queue broker, search server, Node server, external auth, AI provider, cloud storage, or telemetry service. A reverse proxy already operated by the deployer is optional; HTTPS remains required for production. CI inspects the default Compose configuration and performs startup without optional services or outbound SaaS access.

## Dependency growth

At the first working skeleton, commit machine-readable inventories of direct/transitive Go modules and frontend production/development packages. sqlc and Vite are build tools and must still appear in the development inventory. Standard libraries and local packages are not third-party dependencies.

Default permitted unexplained growth per change: zero direct or transitive dependencies. Any increase must have the justification required by AGENTS.md and an updated inventory. CI fails additions absent from the reviewed inventory, lockfile drift, unpinned tool additions, and undeclared production services. Report added/removed counts separately for Go, frontend runtime, and development tooling; do not hide additions behind removals.

An explained addition still must fit size, memory, and service limits. A >10% increase in any size metric versus the accepted main baseline, even below its ceiling, requires recorded review rationale. Do not reset the baseline to conceal regressions. Track critical/high security findings and incompatible licenses as release blockers unless a maintainer records a scoped, time-limited exception.

## CI enforcement and rollout

The skeleton implementation must introduce reproducible check commands and CI jobs for:
1. Foundation-link validation; Go format/vet/tests; Svelte/TypeScript checks; OpenAPI validation once present; deterministic sqlc generation and clean lockfile checks.
2. Production frontend build, manifest/network asset accounting, container build and size, runtime contents, service count, dependency inventory diff, and vulnerability/license review.
3. Real PostgreSQL migration smoke tests, health/readiness checks, then scoped integration and concurrency tests as domain behavior arrives.
4. Controlled-runner startup/RAM measurements, followed by the representative API workload when endpoints exist.

Publish a machine-readable budget report and logs as CI artifacts tied to the exact commit. Required checks fail on ceiling violations or missing required measurements; never substitute zeros. Configure repository required status checks separately when CI is implemented; documenting them does not activate branch protection.

Run deterministic checks on each code change; run affected runtime benchmarks on runtime/query/deployment/dependency changes and the full suite before release. Documentation-only changes require consistency/link review, not fictional runtime measurements. Benchmark regressions >10% against a comparable accepted baseline require investigation and rationale even when below the hard ceiling.

## Exceptions

Prefer removing unnecessary work, shrinking assets, bounding queries, and measuring execution plans before changing budgets. For an exception, record the failing metric, reproducible evidence, simpler alternatives, user/operational impact, proposed limit, responsible owner, and expiry/review date. Obtain explicit maintainer approval before implementation or changing this document. Keep original and new measurements visible. No silent budget increases, disabled checks, or added services to make a benchmark pass.
