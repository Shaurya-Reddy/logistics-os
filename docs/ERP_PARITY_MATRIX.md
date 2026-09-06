# ERPNext functional reference and parity matrix

Status: roadmap/reference document only. It does not expand the current v0.1 implementation scope in PRODUCT.md.

## Purpose

Logistics OS may study ERPNext as a mature functional reference for ERP concepts, workflows, validation rules, edge cases, document relationships, tests, and user terminology. The goal is not to fork ERPNext or reproduce Frappe's framework. The goal is to identify proven business capabilities, keep the useful behavior, simplify the user experience, and reimplement the required behavior independently in Go + Svelte + PostgreSQL.

ERPNext repository reference: https://github.com/frappe/erpnext
Default branch reviewed: develop
License observed on GitHub: GPL-3.0

## Clean-room reference rule

- Do not copy ERPNext source code, DocType JSON, tests, comments, templates, or UI code into this repository by default.
- Concepts, workflows, data relationships, accounting principles, validation behavior, and edge cases may be studied and independently reimplemented.
- When a task references ERPNext, first write the expected behavior in Logistics OS terms before implementation.
- Any proposed direct reuse of GPL-licensed ERPNext code requires explicit maintainer approval and license review before it enters this repository.
- ERPNext and Frappe names/logos are references only and must not be used as Logistics OS branding.

## What the current ERPNext inspection tells us

The current ERPNext tree is organized into broad application areas including accounts, buying, CRM, selling, stock, assets, manufacturing and other domains. The Sales Invoice schema exposes a very large field surface spanning customer identity, dates, POS, returns, accounting dimensions, price lists, stock updates, taxes, discounts, advances, write-offs, loyalty, addresses, commissions, print settings, subscriptions and analytics. The Stock Entry schema combines generic inventory movements with manufacturing, subcontracting, valuation, supplier, printing and accounting concerns.

That breadth is useful as a completeness reference, but it is exactly what Logistics OS should avoid exposing to an ordinary warehouse or billing user. Logistics OS should preserve the underlying capabilities while using task-specific screens, progressive disclosure, and logistics-specific terminology.

## Decision vocabulary

Every ERPNext capability referenced here is classified as one of:

- **KEEP** — same essential business capability is required.
- **SIMPLIFY** — capability is required but the user model/UI should be substantially smaller.
- **SPECIALIZE** — replace the generic ERP concept with a logistics-native workflow.
- **LOCALIZE** — keep out of the country-neutral core; implement in jurisdiction packages.
- **DEFER** — useful later but not needed for the current vertical slice.
- **OMIT** — not part of the intended product unless a real customer requirement appears.

## Phase map

This sequence is intentionally narrower than ERPNext. Later phases do not authorize implementation before PRODUCT.md is updated.

1. **Phase 0 — technical foundation**: application shell, auth, PostgreSQL, migrations, health/readiness, audit plumbing, CI budgets.
2. **Phase 1 — inbound WMS core (current v0.1)**: organization, site, customer, warehouse, location, item, inbound receipt, receive, putaway, stock balances, movement history.
3. **Phase 2 — outbound WMS**: orders, reservation, pick, pack, dispatch, returns, inter-location/site movement.
4. **Phase 3 — contracts and logistics billing**: contracts, rate cards, recurring/activity billing, invoices, credit notes, receivables, statements.
5. **Phase 4 — finance and procurement basics**: chart of accounts, GL, payments, vendor masters, RFQ/PO/receipt/vendor bill.
6. **Phase 5 — CRM/commercial**: lead, opportunity, RFQ, quotation, onboarding, renewal pipeline.
7. **Phase 6 — QMS**: inspection, hold/release, NCR, CAPA, damage/deviation workflows, audit records.
8. **Phase 7 — TMS**: carrier, vehicle, driver, trip, shipment, route, POD, freight verification.
9. **Phase 8 — workforce/assets/admin depth**: shifts, attendance, qualifications, asset maintenance, document templates, richer reporting and localization.

## Parity matrix

| ERPNext concept | Logistics OS concept | Decision | Planned phase | Notes |
|---|---|---|---:|---|
| Company | Organization | KEEP | 1 | Country-neutral company/tenant boundary. |
| Branch / operational unit patterns | Site | SIMPLIFY | 1 | Explicit operational site under organization. |
| Customer | Customer | KEEP | 1 | One customer record shared by CRM, WMS, billing and portal. |
| Supplier | Vendor | KEEP | 4 | Use neutral vendor terminology. |
| Item | Item / SKU / Service | SIMPLIFY | 1/3 | v0.1 stock item first; service items later for billing. |
| UOM | Unit | SIMPLIFY | 1 | One base unit per item in v0.1; conversions deferred. |
| Warehouse | Warehouse | KEEP | 1 | Warehouse belongs to site. |
| Warehouse tree / bin concepts | Location hierarchy | SPECIALIZE | 1 | Receiving, staging, rack/bin, QC hold and dispatch locations. |
| Stock Ledger Entry | Inventory Movement | KEEP | 1 | Immutable movement history; balances are projections. |
| Bin / projected quantity | Stock Balance | KEEP | 1 | Explicit scoped balance table reconciled from movements. |
| Purchase Receipt | Inbound Receipt | SPECIALIZE | 1/4 | Customer-owned inbound first; procurement receipt later. |
| Stock Entry: Material Receipt | Receive | SPECIALIZE | 1 | Operator sees task, not generic stock transaction type. |
| Stock Entry: Material Transfer | Move / Put away | SPECIALIZE | 1/2 | Putaway is a controlled location transfer. |
| Stock Entry: Material Issue | Adjustment / controlled issue | DEFER | 2+ | No generalized adjustment workflow in v0.1. |
| Stock Reconciliation | Inventory Reconciliation | DEFER | 2+ | Requires permissions, reason and compensating movements. |
| Batch | Lot / Batch | DEFER | 2+ | Needed for pharma/agri/expiry-controlled operations. |
| Serial No | Serial Number | DEFER | 2+ | Same principle; scanner-friendly UI. |
| Pick List | Pick Task / Wave | SPECIALIZE | 2 | Warehouse-native picking workflow. |
| Delivery Note | Dispatch | SPECIALIZE | 2 | Operational dispatch record, later linked to transport and billing. |
| Sales Order | Customer Order | SIMPLIFY | 2 | Drives reservation/pick/dispatch when applicable. |
| Sales Invoice | Invoice | SIMPLIFY | 3 | Contract/activity-generated by default; advanced accounting fields hidden. |
| Sales Invoice Return / Credit Note | Credit Note | KEEP | 3 | Linked correcting document; never rewrite posted invoice. |
| Payment Entry | Payment | SIMPLIFY | 3/4 | Simple AR/AP allocation first. |
| General Ledger | General Ledger | KEEP | 4 | Double-entry core; immutable posted entries with reversals. |
| Journal Entry | Journal | SIMPLIFY | 4 | Advanced finance screen, not part of normal warehouse UI. |
| Chart of Accounts | Chart of Accounts | KEEP | 4 | Country-neutral defaults/localized templates later. |
| Cost Center | Cost Centre / Operational Dimension | SIMPLIFY | 4 | Keep only where reporting/accounting requires it. |
| Accounting Dimensions | Dimensions / Tags | DEFER | 4+ | Avoid generic dimensional complexity early. |
| Taxes and Charges | Tax engine | LOCALIZE | 3+ | Neutral tax interface; jurisdiction packages provide rules. |
| Tax Withholding | Withholding engine | LOCALIZE | 3+ | India TDS etc. must not live in core. |
| E-Invoice / statutory documents | Localization adapter | LOCALIZE | 3+ | Jurisdiction-specific integrations only. |
| Quotation | Quote | KEEP | 5 | Connects directly to contract/rate card on win. |
| Lead | Lead | KEEP | 5 | Minimal CRM. |
| Opportunity | Opportunity / RFQ | KEEP | 5 | Logistics requirement capture matters more than generic pipeline fields. |
| Territory / Sales Partner | Commercial attributes | DEFER | 5+ | Only add from real usage. |
| Pricing Rule / Price List | Rate Card | SPECIALIZE | 3 | Logistics units: SFT, pallet, MT, kg, trip, km, man-hour, minimum guarantee, slabs. |
| Subscription / Auto Repeat | Contract Billing Schedule | SPECIALIZE | 3 | Generate monthly/periodic logistics charges from contract rules. |
| Purchase Order | Purchase Order | KEEP | 4 | Standard procurement basics. |
| Material Request | Purchase / Movement Request | SIMPLIFY | 4 | Split by user intent rather than one generic document where useful. |
| Supplier Quotation | Vendor Quote | KEEP | 4 | Needed for procurement comparisons. |
| Purchase Invoice | Vendor Bill | SIMPLIFY | 4 | Standard AP document. |
| Quality Inspection | Inspection | SPECIALIZE | 6 | Embedded in receiving/dispatch rather than isolated paperwork. |
| Quality hold via stock/status workflows | QC Hold / Release | SPECIALIZE | 6 | Inventory state changes automatically with QMS actions. |
| Asset | Asset | KEEP | 8 | Forklifts, scanners, racks, DGs, vehicles and equipment. |
| Asset Maintenance | Maintenance | KEEP | 8 | Preventive/corrective maintenance and downtime. |
| Project | Job / Project | SIMPLIFY | 5+ | Use only when an operation genuinely needs project/job costing. |
| File attachments | Documents / Attachments | KEEP | 1+ | Local durable storage first; optional S3 adapter. |
| Naming Series | Numbering Rules | SIMPLIFY | 1+ | Human pattern configuration such as FSLPL-{YY}-{#####}. |
| Print Formats | Document Templates | SPECIALIZE | 3+ | Block-based, print-safe templates; no HTML editing for ordinary users. |
| Data Import | Import Workbench | SPECIALIZE | 1+ | Saved mappings, required/optional/derived fields, preview and inline errors. |
| Data Export | Export | SIMPLIFY | 1+ | Bounded CSV and later report exports. |
| Role Permission Manager | Permission Bundles | SIMPLIFY | 1 | Deny by default; org/site/customer scope. |
| Workflow | Explicit domain state machines | SPECIALIZE | 1+ | Do not ship a universal low-code workflow engine early. |
| Assignment / ToDo | Task / Exception Queue | SIMPLIFY | 5+ | Command-centre attention model later. |
| Notifications | Notifications / Alerts | DEFER | 5+ | Internal first; external channels optional. |
| Report Builder | Reports | SPECIALIZE | 3+ | Curated operational reports first; generic builder later only if needed. |
| Dashboard / Workspace | Command Centre | SPECIALIZE | 1+ | Tasks, exceptions and KPIs instead of module icon grids. |
| Global Search | Universal Search | KEEP | 2+ | Search customer, item, pallet/serial, invoice, shipment, vehicle, etc. |

## Logistics capabilities that ERPNext does not adequately define for our target product

These are not simple parity items; they are product differentiation and should be designed from logistics operations first:

- Logistics contract as a first-class object linking customer, sites, service scope, SLA, rate card, escalation and billing schedule.
- Activity-based billing from actual warehouse/transport events.
- Storage billing by occupied SFT/SQM, pallet, bin, CBM, weight, day/month and minimum guarantees.
- Shared 3PL warehouse ownership separation across customers.
- Inbound -> inspection -> putaway -> inventory -> outbound -> transport -> POD -> billing as one traceable chain.
- Customer portal showing scoped stock, inbound/outbound, POD, invoices, quality cases and documents.
- Embedded QMS that can hold/release inventory rather than merely attach an inspection record.
- TMS-native shipment, vehicle, driver, carrier, detention, route and POD workflows.
- Warehouse-friendly scanning UX with minimal mandatory fields and fast exception handling.

## ERPNext complexity patterns to avoid

1. **One document carrying unrelated concerns.** Sales Invoice currently spans POS, stock, taxes, advances, loyalty, addresses, commissions, subscriptions and print configuration. Logistics OS should keep the underlying accounting capability but show only the fields relevant to the current task.
2. **Generic stock transaction UI for every intent.** Stock Entry mixes receipt, issue, transfer, manufacture, repack, subcontracting and other purposes. Logistics OS should expose Receive, Put away, Move, Hold, Release, Pick and Dispatch as separate user actions while writing to the same inventory movement core.
3. **Framework-visible implementation details.** Users should not need to understand child tables, DocTypes, posting hooks or internal framework fields.
4. **Country policy in ordinary workflows.** Tax, withholding and statutory rules must remain localization concerns.
5. **Configuration before operation.** The default installation should make the first real warehouse workflow possible with minimal setup.

## Reference backlog

Before implementing each future capability, inspect the corresponding ERPNext behavior and tests, record only the business rules we need, and write Logistics OS acceptance tests from those independently expressed rules.

Priority reference areas:

- Stock ledger, negative-stock prevention, transfers, reconciliation and backdated effects.
- Sales invoice posting, returns/credit notes, outstanding amounts and payment allocation.
- GL posting and reversal behavior.
- Sales order -> delivery -> invoice relationships.
- Purchase order -> receipt -> vendor bill relationships.
- Batch/serial/expiry handling.
- Quality inspection gates and stock status interactions.
- Naming/numbering, CSV import/export and audit history.

When workforce/payroll work begins, review the current Frappe HR/HRMS project separately rather than assuming all HR behavior still lives in ERPNext core.

## Definition of successful parity

Parity does not mean matching ERPNext screen-for-screen or field-for-field. A Logistics OS capability is complete when it supports the required business outcome, preserves accounting/inventory correctness, handles retries/concurrency and auditability, exposes a simpler logistics-native workflow, and remains within the architecture and performance budgets.
