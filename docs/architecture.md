# System Architecture: Schema-Based Multi-Tenancy

## 1. High Level Overview
BizBundl E-Com is a **high-density SaaS platform** designed to serve thousands of e-commerce storefronts from minimal infrastructure.

**Key Pattern:** Schema-Based Multi-Tenancy.
*   **1 Application Binary** (Stateless).
*   **1 Central Database** (PostgreSQL).
*   **N Schemas** (One per Tenant).

---

## 2. Infrastructure Topology

### A. The "Mother" Database Node
A dedicated, high-performance VPS hosting the central data layer.
*   **PostgreSQL**: Holds all tenant data.
    *   `public` schema: Shared data (Users, Billing, System Config).
    *   `shop_123` schema: Isolated tenant data (Products, Orders, Pages).
*   **Redis**: Shared caching layer.
    *   Keys Namespaced: `shop_123:cart:xyz`.
*   **Elasticsearch (Optional)**: Central search cluster (Indices alias per tenant).

### B. The Application Nodes (Stateless)
Multiple cheap VPS nodes running the Go Application.
*   Connected to "Mother Node" via **Private Internal Network** (Low latency, Secure).
*   **Stateless**: No DB stored here. Can be destroyed/recreated instantly.
*   **Scaling**: Add more App Nodes behind the Load Balancer as traffic grows.

---

## 3. Request Lifecyle
1.  **Request**: `shop-abc.bizbundl.com` hits Load Balancer.
2.  **Resolution**: LB routes to `App Node 1`.
3.  **Middleware Identification**:
    *   App parses Host header: `shop-abc`.
    *   Look up `shop_id` from generic cache/DB.
    *   **Context Injection**: `ctx.Set("tenant_id", "shop_123")`.
4.  **Database Connection**:
    *   Middleware executes: `SET search_path TO shop_123, public`.
    *   All subsequent queries (e.g., `SELECT * FROM products`) automatically hit the correct schema.
5.  **Response**: Rentered HTML/JSON returned.

---

## 4. Key Design Decisions

### Why Schema-Based?
*   **Density**: 150 Database-per-tenant = 6GB RAM overhead. 150 Schemas = ~0 RAM overhead.
*   **Cost**: Can host 5,000+ tenants on a single $12 VPS.
*   **Security**: Native Postgres Isolation. Tenant A cannot query Tenant B's tables easily.

### Why Central Database?
*   Enables **Smart Routing / Rolling Migrations**.
*   Allows App Nodes to be treating as "Cattle" (Replaceable).
*   Simplifies Backup (Backup 1 DB instance).

### Why Redis Namespacing?
*   `GET cart:123` becomes `GET shop_abc:cart:123`.
*   Allows shared Redis instance without data collision.
