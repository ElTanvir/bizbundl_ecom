# 4. Multi-Schema Migration Strategy via Worker

Date: 2026-01-08

## Status

Accepted

## Context

Our application uses **Schema-Based Multi-Tenancy** (ADR-0002).
The standard `go-migrate` CLI tool typically connects to a single database schema (usually `public`).
Running `migrate up` against the central database will only migrate the `public` schema (shared tables). It will NOT migrate the hundreds/thousands of isolated tenant schemas (`shop_123`).

We need a reliable way to deploy schema changes to **all** active tenants without manual intervention or downtime.

## Decision

We will implement a dedicated **Migration Worker** (CLI Utility) within the application codebase.

1.  **Discovery**: The worker will query the `public.shops` table (?) or Redis Registry to discover all active Tenant IDs.
2.  **Iteration**: It will iterate through each Tenant ID.
3.  **Execution**: For each tenant:
    *   Construct a dynamic database connection string with `search_path={tenant_id},public`.
    *   Initialize `golang-migrate` instance for that connection.
    *   Run `m.Up()`.
4.  **Concurrency**: For large fleets, we may implement a "Worker Pool" to run migrations in parallel batches (e.g., 10 schemas at a time) to speed up deployment.
5.  **Failures**: A failure in one tenant should be logged and flagged but **must not** stop the deployment for others (unless it's a catastrophic error impacting all).

## Consequences

*   **Automation**: Deployments become "Code Deploy -> Migration Worker -> Done".
*   **Safety**: Ensuring the correct `search_path` is set prevent accidental pollution of the public schema.
*   **Performance**: Migrating 10,000 schemas will take time. We eventually need a "Lazy Migration" strategy (migrate on first access) or massively parallel workers, but for MVP/Growth phase (0-500 shops), a linear/concurrent looper is sufficient.
