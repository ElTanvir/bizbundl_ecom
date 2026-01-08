# Middleware Architecture

BiZBundl uses a "Chain of Responsibility" pattern for middleware.

## 1. Global Middleware (All Requests)
1.  **Recover**: Panics -> 500 Error.
2.  **Logger**: structured logging.
3.  **CORS**: Security headers.
4.  **Compress**: Gzip/Brotli.

## 2. Tenancy Middleware (The Core)
**Goal**: Identify *Who* is calling and *Where* their data is.

**Logic:**
1.  **Extraction**: Parse `Host` header (e.g., `shop-1.bizbundl.com`).
2.  **Lookup**: Check fast cache (Redis) for Tenant Config (ID, Status).
3.  **Context Injection**:
    *   Store `tenant_id` in `c.Locals`.
4.  **Database Configuration**:
    *   This is the critical step for Schema-Based Tenancy.
    *   *Ideally*: The handler uses a wrapper that executes `SET search_path` before queries.
    *   *Or*: The Middleware grabs a connection from the pool, sets the path, and passes it down (Tricky with connection pooling).
    *   **Chosen Approach**: Context Propagation. The Data Access Layer reads `tenant_id` from Context and prepends schema or wraps transaction with `SET search_path`.

## 3. Auth Middleware
1.  **Session Validation**: Paseto/JWT checks.
2.  **RBAC**: Role enforcement.
