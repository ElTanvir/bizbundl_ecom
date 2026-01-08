# 3. Context-Aware Database Transactions

Date: 2026-01-08

## Status

Accepted

## Context

With **Schema-Based Multi-Tenancy** (ADR-0002), strict isolation requires that every database interaction occurs within the correct `search_path`.

PostgreSQL connection poolers (like `pgxpool`) do not guarantee session state (like `search_path`) persists correctly between checkouts in a safe way without explicit reset, and setting it globally on a connection prevents multiplexing.

We need a mechanism to ensure:
1.  Every HTTP request is inextricably linked to a Tenant.
2.  All SQL queries executed during that request use the Tenant's Schema.
3.  We do not accidentally execute queries against `public` (or the wrong tenant) due to connection reuse.

## Decision

We will implement **Context-Aware Database Transactions** via Middleware.

1.  **Middleware Responsibility**: A `TenancyMiddleware` will:
    *   Determine Tenant ID from Hostname.
    *   **Begin a Transaction** immediately.
    *   Execute `SET search_path TO {tenant}, public`.
    *   Inject this **Transaction** (`pgx.Tx`) into the Go `context.Context`.
    *   Commit/Rollback based on HTTP Handler success/failure.

2.  **Store Adaptation**: The Database Store (`internal/store`) will be refactored to:
    *   Accept `context.Context` in all methods.
    *   Check for the presence of a Transaction (`tx`) in the Context.
    *   If present, execute queries via the `tx`.
    *   If absent (e.g., background jobs), fall back to the global Pool (or require a manual Tx).

## Consequences

*   **Safety**: Impossible to "forget" a WHERE clause. The database enforces isolation via the Schema search path.
*   **Consistency**: Atomic requests. If a handler errors, everything is rolled back.
*   **Overhead**: Every request involves a `BEGIN` and `COMMIT`, adding a round-trip (RTT). This is acceptable for the safety guarantees in an e-commerce context.
*   **Refactor Required**: All DB call sites must pass `context.Context` (Standard Go practice anyway).
