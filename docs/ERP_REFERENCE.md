# ERPNext reference and parity plan

Status: canonical reference document. This supplements PRODUCT.md, ARCHITECTURE.md, AGENTS.md, and PERFORMANCE_BUDGET.md. It does not expand approved v0.1 scope.

## Purpose

Logistics OS may study ERPNext as a mature reference for ERP capabilities, workflows, validation rules, document relationships, accounting behavior, inventory edge cases, tests, and terminology. The goal is not to fork ERPNext or reproduce Frappe. The goal is to independently implement the business outcomes we need in Go + Svelte + PostgreSQL while keeping the product logistics-first and lightweight.

Reference repository: https://github.com/frappe/erpnext
Reference branch reviewed: develop
Observed license: GPL-3.0

## Clean-room rules

- Do not copy ERPNext Python, JavaScript, DocType JSON, SQL, tests, fixtures, templates, assets, or documentation prose into this repository by default.
- Use ERPNext to discover capabilities and edge cases, then rewrite the requirement in Logistics OS language before implementation.
- Implement original Go services, SQL, Svelte screens, migrations, and synthetic tests from our own acceptance criteria.
- Do not generate our schema or UI from ERPNext metadata.
- Do not vendor ERPNext, Frappe, Bench, or their runtime.
- Any proposed direct reuse of GPL-licensed source requires explicit maintainer approval and license review before it enters this repository.
- ERPNext and Frappe names/logos are references only and are not Logistics OS branding.

For each future reference exercise, record: upstream URL/revision, observed capability, original Logistics OS requirement, intentional differences, owning module, synthetic acceptance cases, implementation commit, and actual verification.

## What we learned from ERPNext

ERPNext is an excellent completeness reference, but its breadth also demonstrates patterns we explicitly want to avoid.

Sales Invoice spans customer data, dates, POS, returns, accounting dimensions, stock updates, price lists, taxes, discounts, advances, write-offs, loyalty, addresses, commissions, printing, subscriptions and analytics. Stock Entry combines receipt, issue, transfer, manufacturing, repack, subcontracting, valuation and printing concerns.

Logistics OS should preserve needed business capability while exposing smaller task-specific workflows. Ordinary operators should not have to understand generic transaction engines, child-table architecture, posting hooks, framework fields, or country-specific accounting configuration to perform basic logistics work.

## Decision vocabulary

Every ERPNext capability we evaluate is classified as one of:

- **KEEP** — same essential capability is required.
- **SIMPLIFY** — capability is required but the data model or UI should be much smaller.
- **SPECIALIZE** — replace the generic ERP concept with a logistics-native workflow.
- **LOCALIZE** — keep jurisdiction-specific behavior out of the neutral core.
- **DEFER** — useful later, but not part of the current approved slice.
- **OMIT** — not part of the intended product unless a real customer need appears.

## Core parity matrix

| ERPNext concept | Logistics OS concept | Decision | Planned phase | Notes |
|---|---|---|---:|---|
| Company | Organization | KEEP | 1 | Country-neutral tenant/business boundary. |
| Branch / unit patterns | Site | SIMPLIFY | 1 | Explicit operational site. |
| Customer | Customer | KEEP | 1 | Shared by CRM, WMS, billing and portal later. |
| Supplier | Vendor | KEEP | 4 | Neutral terminology. |
| Item | Item / SKU / Service | SIMPLIFY | 1/3 | Stock item first; service item later for billing. |
| UOM | Unit | SIMPLIFY | 1 | One base unit per item in v0.1. |
| Warehouse | Warehouse | KEEP | 1 | Warehouse belongs to Site. |
| Bin / warehouse tree | Location hierarchy | SPECIALIZE | 1 | Receiving, staging, rack/bin, QC hold, dispatch. |
| Stock Ledger Entry | Inventory Movement | KEEP | 1 | Immutable movement history. |
| Bin projected quantity | Stock Balance | KEEP | 1 | Reconciled projection of movement history. |
| Purchase Receipt | Inbound Receipt | SPECIALIZE | 1/4 | Customer custody receipt first; procurement later. |
| Stock Entry: Material Receipt | Receive | SPECIALIZE | 1 | Task-specific operator action. |
| Stock Entry: Material Transfer | Put away / Move | SPECIALIZE | 1/2 | Same inventory core, clearer workflow. |
| Stock Entry: Material Issue | Controlled issue / adjustment | DEFER | 2+ | No generalized adjustment workflow in v0.1. |
| Stock Reconciliation | Inventory Reconciliation | DEFER | 2+ | Permissioned compensating movements only. |
| Batch | Lot / Batch | DEFER | 2+ | Required later for regulated/expiry stock. |
| Serial No | Serial Number | DEFER | 2+ | Scanner-friendly implementation later. |
| Pick List | Pick Task / Wave | SPECIALIZE | 2 | Outbound WMS. |
| Delivery Note | Dispatch | SPECIALIZE | 2 | Operational dispatch record. |
| Sales Order | Customer Order | SIMPLIFY | 2 | Drives reservation/pick/dispatch. |
| Sales Invoice | Invoice | SIMPLIFY | 3 | Prefer contract/activity-generated billing. |
| Credit Note / Return | Credit Note | KEEP | 3 | Linked correction, never rewrite posted invoice. |
| Payment Entry | Payment | SIMPLIFY | 3/4 | Basic AR/AP allocation first. |
| General Ledger | General Ledger | KEEP | 4 | Double-entry, immutable posted entries. |
| Journal Entry | Journal | SIMPLIFY | 4 | Advanced finance screen. |
| Chart of Accounts | Chart of Accounts | KEEP | 4 | Neutral defaults; localized templates later. |
| Cost Center | Cost Centre / Operational Dimension | SIMPLIFY | 4 | Only where genuinely required. |
| Accounting Dimensions | Dimensions / Tags | DEFER | 4+ | Avoid early generic complexity. |
| Taxes and Charges | Tax Engine | LOCALIZE | 3+ | Jurisdiction packages supply rules. |
| Tax Withholding | Withholding Engine | LOCALIZE | 3+ | India TDS etc. stays outside core. |
| E-Invoice / statutory docs | Localization Adapter | LOCALIZE | 3+ | Country-specific integrations. |
| Lead | Lead | KEEP | 5 | Minimal CRM. |
| Opportunity | Opportunity / RFQ | KEEP | 5 | Capture logistics requirements. |
| Quotation | Quote | KEEP | 5 | Converts to contract/rate card on win. |
| Pricing Rule / Price List | Rate Card | SPECIALIZE | 3 | SFT, pallet, MT, kg, trip, km, man-hour, slabs, minimum guarantee. |
| Subscription / Auto Repeat | Contract Billing Schedule | SPECIALIZE | 3 | Periodic logistics billing. |
| Purchase Order | Purchase Order | KEEP | 4 | Procurement basics. |
| Supplier Quotation | Vendor Quote | KEEP | 4 | Procurement comparison. |
| Purchase Invoice | Vendor Bill | SIMPLIFY | 4 | Standard AP document. |
| Quality Inspection | Inspection | SPECIALIZE | 6 | Embedded in receiving/dispatch. |
| Quality hold patterns | QC Hold / Release | SPECIALIZE | 6 | Inventory availability enforced by QMS. |
| Asset | Asset | KEEP | 8 | Forklifts, scanners, racks, DGs, vehicles. |
| Asset Maintenance | Maintenance | KEEP | 8 | Preventive/corrective maintenance. |
| Project | Job / Project | SIMPLIFY | 5+ | Only if operationally useful. |
| File Attachments | Documents / Attachments | KEEP | 1+ | Local durable storage first. |
| Naming Series | Numbering Rules | SIMPLIFY | 1+ | Example: FSLPL-{YY}-{#####}. |
| Print Formats | Document Templates | SPECIALIZE | 3+ | Block-based, print-safe editor; no routine raw HTML. |
| Data Import | Import Workbench | SPECIALIZE | 1+ | Saved mappings, preview, inline errors. |
| Data Export | Export | SIMPLIFY | 1+ | Bounded CSV/report exports. |
| Role Permission Manager | Permission Bundles | SIMPLIFY | 1 | Deny by default, scoped by org/site/customer. |
| Workflow | Explicit Domain State Machines | SPECIALIZE | 1+ | No general low-code workflow engine early. |
| Assignment / ToDo | Task / Exception Queue | SIMPLIFY | 5+ | Command-centre model later. |
| Notifications | Alerts | DEFER | 5+ | Internal first; external channels optional. |
| Report Builder | Reports | SPECIALIZE | 3+ | Curated operational reports first. |
| Dashboard / Workspace | Command Centre | SPECIALIZE | 1+ | Tasks, exceptions and KPIs, not icon grids. |
| Global Search | Universal Search | KEEP | 2+ | Customer, item, pallet/serial, invoice, shipment, vehicle, etc. |

## Logistics-native capabilities beyond simple ERP parity

These are core differentiators and should be designed from logistics operations first:

- Logistics Contract as a first-class object linking customer, sites, service scope, SLA, rate card, escalation and billing schedule.
- Activity-based billing from actual warehouse and transport events.
- Storage billing by SFT/SQM, pallet, bin, CBM, weight, day/month and minimum guarantee.
- Shared 3PL warehouse inventory separation by customer ownership.
- One traceable chain: inbound -> inspection -> putaway -> inventory -> outbound -> transport -> POD -> billing.
- Customer portal for scoped stock, inbound/outbound, POD, invoices, quality cases and documents.
- Embedded QMS that can hold/release stock availability.
- TMS-native shipment, carrier, vehicle, driver, detention, route and POD workflows.
- Warehouse-scanner UX with minimal fields and fast exception handling.

## Complexity patterns to avoid

1. **One document carrying unrelated concerns.** Keep the accounting capability, but show only fields relevant to the task.
2. **One generic stock transaction screen for every intent.** Expose Receive, Put away, Move, Hold, Release, Pick and Dispatch as clear actions over one movement core.
3. **Framework-visible implementation details.** Users should not need to understand DocTypes, child tables, hooks, schema metadata, or internal IDs.
4. **Country-specific rules inside ordinary workflows.** Taxes, withholding and statutory behavior belong in localization modules.
5. **Configuration before operation.** Default setup should reach a real warehouse workflow quickly.

## Ordered implementation roadmap

Later phases are planning references only until PRODUCT.md explicitly expands scope.

### S0 — Running technical skeleton

Build only:
- Go executable and HTTP server.
- Svelte + TypeScript + Vite shell, compiled to static assets served by Go.
- PostgreSQL connection through pgx.
- Explicit versioned SQL migrations and migration command.
- `/health` and DB/schema-aware `/ready`.
- Docker Compose with exactly application + PostgreSQL.
- OpenAPI starter contract.
- Initial CI, dependency inventory and performance-budget checks.

Exit criteria:
- `docker compose up` starts the two required services.
- Browser loads Logistics OS shell.
- `/health` returns process liveness without DB dependency.
- `/ready` fails until PostgreSQL/schema is ready, then succeeds.
- Runtime image contains no Node runtime or frontend source tree.
- Initial size/service checks produce real measurements rather than placeholders.

### S1 — Identity and scoped master data

After S0:
- Local auth/bootstrap and permission grants.
- Organization -> Site -> Customer -> Warehouse -> Location -> Item.
- OpenAPI, migrations, scoped SQL and minimal screens.

Exit criteria include denied cross-organization/site access, stale update rejection, and safe disablement of referenced records.

### S2 — Draft/open inbound receipts and Receive

After S1:
- Draft/open receipt lifecycle.
- Exact-decimal positive quantities.
- Partial receiving into staging.
- Idempotent posting, immutable movements, audit and concurrent-update protection.

Exit: retries do not duplicate stock, over-receipt is rejected, invalid transitions are blocked, and concurrent receiving preserves correct quantity.

### S3 — Putaway, inquiry and correction

After S2:
- Putaway from staging to storage.
- Stock balances and movement-history inquiry.
- Bounded CSV export.
- Permissioned linked reversals and short-receipt closure.

Exit: no negative stock, balances reconcile from movements, competing putaways are safe, and corrections never overwrite posted history.

### S4 — v0.1 qualification

Complete access, concurrency, precision, restore and performance evidence tied to the tested commit. Resolve failures without weakening budgets.

After S4, choose one next operator outcome—most likely outbound WMS or contract/service billing—and update PRODUCT.md before implementation.

## Original inbound acceptance cases

Use synthetic customer C1, item I1, base unit EACH, receiving location R and storage location A in one authorized warehouse.

- Receipt expects 10.000000; receive 6.000000: staging becomes 6 and receipt is partially received.
- Retry the same request with the same idempotency key: result is replayed and stock remains 6. Reuse the key with a different payload: reject.
- Put away 4: staging becomes 2, storage becomes 4, total remains 6.
- Two competing putaways of 4 from staging 6 cannot both succeed or produce a negative balance.
- Receive the remaining 4: accepted total becomes 10. A further receipt of 1 is rejected.
- A short receipt may close only with appropriate permission and a recorded reason after accepted stock is fully put away.
- Reject 0, negative values and precision beyond 6 decimal places instead of silently rounding.
- A correction creates linked compensating movements and a reason; duplicate or excess reversal cannot remove the same stock twice.
- A second organization or unauthorized site cannot read or mutate these records. Customer C2 stock cannot satisfy C1 references.

These are original Logistics OS acceptance requirements, not copied ERPNext tests.

## Reference backlog for later phases

Before implementing each later capability, inspect the corresponding ERPNext behavior and tests and record only the business rules we need.

Priority reference areas:
- Stock ledger, negative-stock prevention, transfers, reconciliation and backdated effects.
- Sales invoice posting, returns/credit notes, outstanding amounts and payment allocation.
- General ledger posting and reversal behavior.
- Sales order -> delivery -> invoice lineage.
- Purchase order -> receipt -> vendor bill lineage.
- Batch/serial/expiry handling.
- Quality inspection gates and stock-status interactions.
- Naming/numbering, CSV import/export and audit history.

When workforce/payroll work begins, evaluate the current Frappe HR/HRMS project separately rather than assuming HR behavior remains in ERPNext core.

## Definition of successful parity

Parity does not mean matching ERPNext screen-for-screen or field-for-field. A Logistics OS capability is complete when it delivers the required business outcome, preserves inventory/accounting correctness, handles retries/concurrency and auditability, exposes a simpler logistics-native workflow, and remains within ARCHITECTURE.md and PERFORMANCE_BUDGET.md.
