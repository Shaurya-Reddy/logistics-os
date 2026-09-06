# Product foundation

Status: agreed direction and v0.1 scope; this document does not claim implementation exists.
Read with [architecture](ARCHITECTURE.md), [agent rules](AGENTS.md), and [performance budgets](PERFORMANCE_BUDGET.md).

## Purpose and intended users

Logistics OS is a lightweight, open-source, self-hosted operating system for logistics operations. It serves small and medium warehouse operators, third-party logistics providers, and businesses operating their own warehouses. Operators own their data and can run the core without a cloud subscription, vendor account, or AI provider.

The long-term scope connects ERP (operational and financial records), WMS (warehouse and inventory), TMS (transport), QMS (quality), CRM (customer relationships), and workforce operations. These are a roadmap, not permission to build six suites at once. Open-source licensing must be selected by the maintainer and recorded in LICENSE before a public release; this document alone grants no license.

## Product principles

- Make a small operation productive on modest hardware with simple installation, backup, and upgrades.
- Prefer one complete, reliable workflow over many incomplete modules.
- Keep inventory quantities, ownership, permissions, and auditability correct under concurrent use.
- Use a country-neutral core; local tax, statutory, address, and document rules belong in explicit localization modules.
- Keep core workflows usable without internet access from the server. A browser connection to the server is required; offline synchronization is not promised.
- Offer documented APIs and portable data; avoid vendor lock-in and mandatory external services.
- AI is optional, disabled by default, and assistive. It cannot be required to operate, authorize transactions, or silently change stock.
- Treat speed, accessibility, and operational clarity as product requirements.

## First vertical workflow: receive customer-owned stock into a warehouse

The v0.1 user outcome is: an authorized operator can receive a customer's goods, put them away, and explain the resulting stock balance from an immutable movement history.

Build in bounded milestones:
1. Technical skeleton: Go server, Svelte/TypeScript/Vite shell, PostgreSQL connection, versioned SQL migrations, health/readiness endpoints, embedded production assets, Docker Compose, and initial budget checks. No business modules in this task.
2. Master data: Organization -> Site -> Customer -> Warehouse as the implementation sequence. A site belongs to an organization; a warehouse belongs to a site; a customer belongs to an organization and may hold stock in multiple warehouses. Customers do not own the warehouse hierarchy.
3. Inbound slice: customer-owned items/SKUs and units, warehouse locations including receiving/staging, draft inbound receipt and lines, partial receiving, putaway, stock inquiry, and audit history.

Receipt lifecycle: draft -> open -> partially received -> received -> closed. An open receipt with no movements may be cancelled. Receipt closure requires all accepted units to be put away and any short receipt to have a recorded reason. Posting receipt quantities creates stock in staging; putaway transfers it to a storage location without changing total stock. Repeated submissions must not duplicate movements. Over-receipt is rejected in v0.1. Posted mistakes require a permission-controlled compensating movement with a reason; never overwrite history or create negative stock. No generalized adjustment workflow is implied.

The minimum stock identity is organization, customer, warehouse, location, item, and unit. v0.1 uses one base unit per item; fractional quantities use exact decimals. Lot, serial, expiry, and unit conversion features are deferred. The UI must make this limitation clear for goods needing those controls.

## v0.1 acceptance

- An administrator configures organization/site/customer/warehouse/location/item records and assigns users to permitted sites.
- An inbound operator creates and receives a receipt; a warehouse operator puts stock away; a supervisor reviews balances and movement history.
- Partial receipts and short closure are explicit. Concurrent requests and retries preserve quantities; cross-customer and cross-organization access is denied.
- Authorized users can search and paginate stock and receipts, view actors/timestamps/reasons, and export a bounded CSV stock list.
- A fresh installation, database migration, backup/restore exercise, and the performance checks succeed without an external SaaS dependency.
- Tests cover the main workflow, denied access, invalid transitions, reversal behavior, retries, and concurrent stock updates.

## Roles

Permissions are enforced on the server and scoped by organization and, where applicable, site. Roles are permission bundles, not hard-coded checks against role names. Users may hold more than one bundle.
- Organization administrator: manage organization configuration, master data, users, and grants within that organization.
- Operations supervisor: oversee inbound work, close short receipts, authorize reversals, and inspect operational audit records at assigned sites.
- Inbound operator: maintain drafts and receive goods at assigned sites; cannot grant access or authorize reversals.
- Warehouse operator: perform putaway and view relevant stock at assigned sites.
- Auditor/viewer: read scoped stock, documents, movement history, and approved exports; no writes.

Deployment administration is an infrastructure responsibility and does not automatically confer unrestricted business access. Customer portal, driver, accountant, and HR roles are future scope.

## Explicitly out of scope for v0.1

Outbound orders, reservations, picking/packing/shipping, inter-site transfers, procurement, billing, general ledger, payroll, route optimization, fleet telematics, customs, statutory reporting, tax engines, advanced CRM, QMS workflows, workforce scheduling, customer portals, EDI, real-time integrations, lot/serial/expiry tracking, and general-purpose workflow builders are excluded.

Also excluded: mandatory AI, autonomous stock decisions, native mobile apps, offline-first sync, microservices, an extension marketplace, multi-region operation, and enterprise-scale availability promises. Organization isolation is required, but a public multi-tenant SaaS control plane is not.

## UX philosophy

Design around warehouse work: clear task queues, searchable compact lists, persistent organization/site context, visible status, quantities, units, and customer ownership. Use familiar words such as Receive and Put away. Make the next valid action obvious and explain blocked actions.

Support keyboard navigation, visible focus, labeled forms, accessible contrast, and errors next to the affected input. Never rely on color alone. Barcode scanner keyboard input may reuse ordinary item fields; specialized hardware integration is deferred.

Use responsive layouts for desktops and warehouse tablets. Preserve entered data after errors, distinguish loading/empty/error states, prevent duplicate submissions, and request confirmation for consequential actions. Display local dates and time zones without obscuring the underlying event time. Avoid dashboard decoration, giant component suites, and endless configuration before the first receipt.

## Delivery discipline

Finish and measure the skeleton before domain work. Deliver each milestone in independently reviewable tasks with demonstrated acceptance criteria. Broaden scope only through an explicit product decision and corresponding updates to these foundation documents.
