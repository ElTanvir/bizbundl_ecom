# Database Migrations

**Philosophy**: "Additive Changes Only".

Even though our [Deployment Strategy](deployment.md) allows for breaking changes via Version Routing, we should still strive for **Additive Migrations** as a best practice. It reduces risk if something goes wrong with routing.

## The "Additive" Rule
**Never break the potential for the OLD code to run, if possible.**

1.  **Don't Rename Columns**: Add a new column instead.
2.  **Don't Delete Columns**: Mark them deprecated in comments, delete in a future cleanup release.
3.  **Don't Add `NOT NULL` constraints** to existing tables without a default value.

## Handling Multi-Schema Migrations
Since we have N schemas (`shop_1`, `shop_2`...):

1.  **Do NOT** run `migrate up` on the DB globally found on the internet.
2.  **The Migration Worker** (See [ADR 0004](adr/0004-multi-schema-migration-strategy.md)):
    *   We use a dedicated CLI tool: `cmd/migrate_worker`.
    *   It iterates over all Tenants in the Registry.
    *   For each tenant, it sets `search_path` and runs the migration.
    *   **Usage**: `go run cmd/migrate_worker/main.go`

## Startup Migration (Development)
In Development/Local:
*   The App automatically migrates "known static tenants" on startup for convenience.

## Lazy Migration (Production)
In Production:
*   Use the **Smart Router** approach.
*   Migrate batches of tenants *before* routing them to the new code version.
