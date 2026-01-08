# Forking Strategy: Adapting for Single Vendor

**Goal:** Transform the BizBundl Multi-Tenant SaaS Engine into a Standalone Single-Vendor Store.

## 1. Why Fork?
You might want to deploy this engine for a specific high-volume client (Enterprise) who wants their own dedicated infrastructure/codebase, or you might want to pivot the business model.

## 2. The Core Difference
*   **SaaS (Current)**:
    *   Identify Tenant via `Host` Header.
    *   Dynamic Schema Switching (`SET search_path`).
    *   Strict Data Isolation.
*   **Single Vendor (Fork)**:
    *   Tenant Identity is static (Always "Default Shop").
    *   Schema is static (`public` or specific schema).

## 3. Adaptation Steps (The 5-Minute Refactor)

You do **NOT** need to rewrite the codebase. The "Schema-Based" architecture makes this trivial.

### Step A: Middleware Override
Locate `internal/middleware/tenancy.go`.
*   **Current Logic:** `tenantID := ExtractFromHost(c.Hostname())`
*   **Change To:** `tenantID := "default"`

```go
// Adaptation
func TenancyMiddleware(c *fiber.Ctx) error {
    // HARDCODE THE TENANT
    c.Locals("tenant_id", "public") // Or specific schema name
    return c.Next()
}
```

### Step B: Database Connection
*   The system currently runs `SET search_path TO {tenant_id}`.
*   If you hardcode `tenant_id = "public"`, it runs `SET search_path TO public`.
*   **Result:** It behaves exactly like a normal, single-tenant Application.

### Step C: Configuration
*   Update `config.yaml` or Env Vars to remove "Multi-Tenant" flags if you add any specific ones.
*   Set `REDIS_PREFIX` to empty string `""` if you don't wantnamespacing, or keep it `shop:` to be safe.

## 4. Maintenance Benefit
Because the change is so minimal (1 line in Middleware), you can potentially **keep the codebase synced** with the SaaS version.
*   Keep the logic allowing dynamic tenancy.
*   Just use using an Environment Variable `FORCE_TENANT_ID=my_shop`.
*   If set, Middleware ignores Host header and uses Env Var.
*   **Result**: One codebase serves both SaaS and Single Vendor clients.
