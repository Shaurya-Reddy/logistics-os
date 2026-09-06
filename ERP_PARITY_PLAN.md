# ERP parity and clean-room reference plan

Status: documentation proposal, 2026-09-06. No ERP feature is implemented by this document.
This plan supplements [PRODUCT.md](PRODUCT.md), [ARCHITECTURE.md](ARCHITECTURE.md), [AGENTS.md](AGENTS.md), and [PERFORMANCE_BUDGET.md](PERFORMANCE_BUDGET.md); it does not expand approved v0.1 scope or override those files.

## Evidence and current baseline

Inspected logistics-os main at [99ef659](https://github.com/Shaurya-Reddy/logistics-os/commit/99ef65998afbdbd3b111fb6ee03d3db327f5c5ea). Its complete recursive tree contains five files: the four populated foundation documents and a two-line README. There are no application sources, Go/frontend manifests, migrations, OpenAPI contract, tests, containers, or CI workflows. All functional areas below are unimplemented. Earlier conversation reports of empty foundations are superseded by this snapshot.

- PRODUCT.md already establishes customer-owned inbound stock, partial receipt, putaway, scoped access, stock inquiry, CSV export, and audit. It deliberately excludes procurement, sales fulfillment, finance, QMS, TMS, and workforce in v0.1. ERP breadth is a long-term direction, not a missing first-release requirement.
- ARCHITECTURE.md already defines module ownership, exact quantities, transactional movement posting, reversals, idempotency, auth, and two-service deployment. It has no implemented schema. Future financial ownership, valuation, and commercial document contracts remain undecided.
- AGENTS.md already limits changes to bounded slices, requires dependency justification and verification, and reserves scope/architecture exceptions for maintainer decisions. This plan adds a reference discipline without changing those rules.
- PERFORMANCE_BUDGET.md already specifies measurable limits and enforcement requirements. There are no measured results or enforcement jobs yet. Adding this plan does not establish a benchmark pass.
- README.md identifies the project but does not provide installation instructions; installation documentation should accompany the working skeleton.
- There is no LICENSE file in this snapshot. The maintainer must select the project license before public release; this plan selects none.

ERPNext reference snapshot: [c1e0865](https://github.com/frappe/erpnext/commit/c1e08657710d24d221675de2a523864686545bcc), the default development head observed on 2026-09-06, not a stable-release compatibility target. Inspected its module inventory and module/DocType directory names, plus public product and stock documentation. No controllers, algorithms, schema bodies, fixtures, or UI implementation were imported. Directory presence establishes a reference area, not verified behavioral equivalence.

## Reference boundary

Use ERPNext to discover business capabilities and questions. Write original requirements, Go services, SQL, Svelte screens, and synthetic tests from Logistics OS use cases. Do not copy or translate ERPNext Python/JavaScript, DocType JSON, SQL, tests, fixtures, templates, assets, or documentation prose. Do not generate our schema or UI from upstream metadata. Do not vendor ERPNext or adopt Frappe/Bench as a runtime.

ERPNext's [license file](https://github.com/frappe/erpnext/blob/c1e08657710d24d221675de2a523864686545bcc/license.txt) identifies GPLv3. “Clean-room” here describes a no-copy engineering process, not a legal certification or a claim of independently isolated teams. This task inspected public repository metadata and documentation. Any proposed source reuse must be handled separately with maintainer review; reference links confer no license on this project.

For every future slice, record: reference URL and revision/access date; observed capability; original requirement; intentional differences; affected owner; synthetic acceptance cases; implementation commit and actual verification. Clearly label behavior inferred from names until confirmed through documentation or permitted observation. Keep upstream snippets and copied test data out of implementation prompts and commits.

## Functional coverage and deliberate differences

The [module inventory](https://github.com/frappe/erpnext/blob/c1e08657710d24d221675de2a523864686545bcc/erpnext/modules.txt) is the breadth checklist. “Later” below means a proposal requiring a scope decision, not an approved delivery commitment.

### Setup, identity, and shared records — v0.1 subset

Reference: [Setup](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/setup), [Customer](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/selling/doctype/customer), [Item and Warehouse](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/stock/doctype).
Build Organization, Site, Customer, Warehouse, Location, customer-owned Item, and one base unit per item. Identity owns users/sessions/grants; Organization, Customers, and Warehouse retain their existing ownership. Use stable IDs and scoped codes; do not recreate a generic DocType engine.
Exit: a user can configure the scoped hierarchy, and both API and database constraints reject cross-organization references. Broader configuration, custom fields, portal access, and workflow builders remain deferred.

### Stock and warehouse management — first functional priority

Reference: [Stock directory](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/stock/doctype), including purchase_receipt, putaway_rule, stock_entry, stock_ledger_entry, bin, stock_reservation_entry, pick_list, delivery_note, batch, and serial_no; [stock guide](https://docs.frappe.io/erpnext/stock).
Build only inbound receipts, staging, manual putaway, balances and movement history in v0.1. Warehouse orchestrates; Inventory owns movement/balance writes.
Intentional difference: customer goods entering a 3PL warehouse are custody receipts, not necessarily purchases by the operator. Receipt posting must not create supplier debt or treat customer-owned stock as an operator asset. Keep quantity custody distinct from future financial valuation.
Later candidates: outbound reservation/pick/pack/dispatch, returns, stock counts, lot/serial/expiry, and unit conversion. Each needs its own product decision and invariants.

### Selling and CRM — customer master now; commercial workflow later

Reference: [Selling](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/selling/doctype) and [CRM](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/crm).
Customer master is in v0.1. Leads, opportunities, quotations, sales orders, pricing, and credit control are not.
Candidate outcome: convert an accepted logistics-service quotation into an order with explicit billing terms. Do not confuse selling warehousing services with selling the customer's goods. Exit for a future slice: preserve quotation/order lineage and authorized state changes; independently test pricing and rounding. Sales invoicing belongs to future finance design.

### Buying — later

Reference: [Buying](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/buying/doctype), including supplier, request_for_quotation, supplier_quotation, and purchase_order; purchase_receipt is under Stock, purchase_invoice under Accounts.
Candidate outcome: procure operator-owned supplies or services with supplier/order/receipt/invoice traceability. Define a procurement owner before implementation; do not assign procurement tables to Warehouse by analogy.
Exit: distinguish partial delivery, outstanding order quantities, and supplier invoice matching. A customer custody receipt must never satisfy an operator purchase order implicitly.

### Accounts and finance — later, separate design gate

Reference: [Accounts](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/accounts/doctype), including account, journal_entry, gl_entry, sales_invoice, purchase_invoice, payment_entry, fiscal_year, cost_center, bank_reconciliation_tool, and budget.
Candidate sequence: agree service billing rules and currencies; design balanced journals and periods; then implement one invoice/payment flow. Finance owns monetary postings; warehouse history remains operational evidence.
Require explicit decimal precision, currency rounding, account mapping, posting dates, period locks, corrections, and idempotency. Acceptance must include balanced debits/credits, partial payments, reversals, and denied closed-period posting. Tax/statutory behavior stays in localization. Do not claim accounting parity from an invoice CRUD screen.

### Quality, support, and maintenance — later

Reference: [Quality Management](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/quality_management/doctype), [Stock quality_inspection](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/stock/doctype/quality_inspection), [Support](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/support), and [Maintenance](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/maintenance).
Candidates: inspection results and non-conformance/corrective actions linked to operational records, service issues, and equipment maintenance. Define whether a quality hold affects availability before introducing it; Inventory must enforce any resulting stock restriction.
Exit: trace an inspection failure to its resolution with permissions and history. No quality hold, support workflow, or maintenance scheduler is implied in v0.1.

### Projects and assets — later, lower priority

Reference: [Projects](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/projects) and [Assets](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/assets).
Evaluate contract/project cost attribution and operator-owned equipment registers after operational and finance foundations. Asset depreciation requires finance policy; timesheets do not establish payroll support. Exit criteria must be agreed for a bounded feature before implementation.

### Manufacturing and subcontracting — outside present logistics roadmap

Reference: [Manufacturing](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/manufacturing) and [Subcontracting](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/subcontracting).
Do not build production planning, bills of materials, or work orders merely to match ERPNext breadth. Future warehouse kitting/value-added services need a separate proposal distinguishing custody transformations from owned production and financial costing.

### Regional, portal, integration, and supporting modules — deferred

Reference: [Regional](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/regional), [Portal](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/portal), [ERPNext Integrations](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/erpnext_integrations), [EDI](https://github.com/frappe/erpnext/tree/c1e08657710d24d221675de2a523864686545bcc/erpnext/edi), and the module inventory's Utilities, Communication, Telephony, and Bulk Transaction areas.
Use country-neutral core records and future explicit localization adapters. External messages/imports need bounded jobs, authorization, idempotency, and auditable failures. Portal access requires a customer-scoping design. Use PostgreSQL jobs inside the same app when needed; do not add brokers or unused integration scaffolding.

### TMS and workforce — Logistics OS extensions, not established ERPNext parity

Stock includes delivery_trip and shipment reference areas, but this inspection does not establish a complete transport suite. Driver assignment, proof of delivery, route optimization, and telematics require separate discovery.
The inspected ERPNext module inventory does not include HR/payroll. Do not infer workforce coverage from the brand or historical ERPNext versions; evaluate the relevant separate product/repository if that work is requested. Both areas remain outside v0.1.

## Ordered implementation backlog

All tasks below are pending. Finish each gate before its dependents; do not implement this entire list in one task.

1. **S0 — running skeleton.** Go entry point, static Svelte/TypeScript/Vite shell served by Go, pgx connection, explicit migrations subcommand, /health and database/schema-aware /ready, and two-service Compose. Pin and justify tool/dependency versions. Add relevant CI and dependency/budget reports under PERFORMANCE_BUDGET.md. Demonstrate fresh startup, migration failure, runtime image contents, and size measurements. No domain modules or speculative jobs.
2. **S1 — identity and scoped master data.** After S0, introduce local auth/bootstrap and grants as prerequisites to exposed business writes; then Organization -> Site -> Customer -> Warehouse, followed by locations/items for inbound. Preserve existing module ownership. Add OpenAPI, scoped SQL, constraints, and minimal screens. Demonstrate denied cross-organization/site access, stale edits, and disabled referenced records.
3. **S2 — draft/open receipts and receive.** After S1, implement the PRODUCT.md lifecycle and positive exact-decimal lines. Post partial quantities to staging in one transaction with Inventory, audit, and idempotency records. Demonstrate retries, over-receipt rejection, invalid transitions, and concurrent receiving.
4. **S3 — putaway, inquiry, and correction.** After S2, transfer staging quantities atomically, add bounded stock/history inquiry and CSV export, linked supervisor-authorized reversals, and short-receipt closure. Reconcile balances from movements and prove no negative stock.
5. **S4 — v0.1 qualification.** Complete access/concurrency/precision coverage, restore exercise, and required budget measurements with evidence tied to the tested commit. Resolve failures without weakening budgets.
6. **Future decision — outbound or service billing.** Only after v0.1 qualification, select the next operator outcome and update foundations explicitly. The broad ERP areas above are discovery input, not permission to expand S0–S4.

## Original inbound acceptance examples

Use synthetic customer C1, item I1, base unit EACH, receiving location R, and storage location A within one authorized warehouse.

- Receipt expects 10.000000; receive 6.000000: staging is 6, history has one receipt posting, and document is partially received.
- Retry the identical receive request with its idempotency key: same result, still 6 in staging. Reuse that key with another payload: reject.
- Put away 4: staging 2, storage 4, total 6. Competing putaways of 4 from staging 6 must not both succeed or leave a negative balance.
- Receive the remaining 4: accepted total 10; put away the remaining staging stock before closure. A further receipt of 1 is rejected.
- On a separate receipt expecting 10 with only 6 accepted and fully put away, closure requires the appropriate permission and a short-receipt reason.
- Reject quantity 0, negative values, and 0.0000001 rather than silently rounding.
- A permitted correction creates linked compensating movements and a reason; replay or excess reversal cannot remove the same stock twice. Follow location availability when goods have moved.
- A second organization or unauthorized site cannot read or mutate these records. Customer C2 stock cannot be substituted for C1 in receipt lines or movement references.

These cases are newly authored requirements from the existing foundations, not copied ERPNext tests.

## Architecture and verification guardrails

Preserve Go + Svelte/TypeScript/Vite + PostgreSQL, pgx/sqlc, versioned REST/OpenAPI, explicit module services, and one application runtime. Node remains build-only. No ORM, Frappe runtime, Redis, Kafka, RabbitMQ, Elasticsearch, Kubernetes, or microservice split.

Unchanged budget highlights: initial JS <=200 KB gzip, initial CSS <=50 KB gzip, runtime image <=100 MB uncompressed, idle application RSS <=150 MB, startup p95 <2 seconds under the specified conditions, exactly two mandatory running services. The full measurement contract, total asset limits, and endpoint latency limits in PERFORMANCE_BUDGET.md remain authoritative.

For this documentation change: reread all foundations, verify source references and relative links, and confirm the diff adds only this plan. Reopen the committed file to verify persistence. Application tests and runtime measurements are not applicable because no application exists. For future implementation, attach actual checks and measurements; never change a reference checklist item to “implemented” without evidence.
