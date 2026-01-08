# 2. Schema-Based Multi-Tenancy

Date: 2026-01-08

## Status

Accepted

## Context

We are building a multi-tenant e-commerce platform (`bizbundl_ecom`). We need a strategy to isolate tenant data while maintaining high density, low operational cost, and ease of maintenance for a solo developer or small team.

Options considered:
1.  **Database per Tenant**: Highest isolation, but high resource overhead (memory/connections) and complex operations (backups/migrations for 1000s of DBs).
2.  **Row-level Security (RLS) / Discriminator Column**: Lowest overhead, but weak isolation (easy to leak data via missed `WHERE` clause) and harder to backup/restore single tenants.
3.  **Schema per Tenant**: Balanced approach. Identifying tenants via PostgreSQL `search_path`.

## Decision

We will use **Schema-Based Multi-Tenancy** on a **Central PostgreSQL Database**.

*   Each shop gets its own schema: `shop_{id}`.
*   Shared data (users, platform config) lives in `public`.
*   We will use a "Mother Node" (Central DB) topology initially.
*   Capacity scaling will be handled by "Registry-Based Routing" to secondary DB nodes (sharding at the application level) when the primary conceptual node is full.

## Consequences

*   **Pros**:
    *   **High Density**: Can host thousands of inactive shops with effectively zero overhead.
    *   **Strong Isolation**: `SET search_path` guarantees queries typically default to the tenant's schema, reducing data leak risks compared to RLS.
    *   **Simplicity**: Backups can be per-schema or full DB. Migration tools (Golang-migrate) support schema iterations.
*   **Cons**:
    *   **Migration Complexity**: Migrations must run across thousands of schemas. This requires robust tooling.
    *   **Connection Pooling**: Managing connection pools effectively when switching search paths requires care (handled via Transaction Wrapping).
