# Deployment: Zero-Downtime Smart Routing

**Goal:** Updates without maintaining "Backward Compatibility" code logic.

## 1. The Strategy: "Schema-Aware Routing"
Instead of "All Traffic goes to New Version", we use a **Smart Load Balancer** (or Nginx Map) to route specific tenants to specific App Versions.

Because we have a **Central Database** with isolated schemas, multiple versions of the App can connect to the *same* database simultaneously, as long as they talk to *different schemas*.

## 2. The Upgrade Workflow (Rolling)

Let's say we are upgrading from `v1.0` to `v1.1`.

### Step 1: Provision
*   **Cluster A (Old)**: Running `v1.0`. Serving all tenants (Schemas v1).
*   **Cluster B (New)**: Start up `v1.1`. Connects to Central DB. Serving 0 tenants.

### Step 2: Migrate & Switch (Per Tenant or Batch)
Pick a batch of 10 shops (e.g., `shop_X`, `shop_Y`).

1.  **Migrate Schemas**:
    *   Run migration script ONLY on schemas `shop_X`, `shop_Y`.
    *   Upgrade them from `schema_v1` -> `schema_v2`.
    *   *Note: `shop_Z` (still on Cluster A) is unaffected.*

2.  **Switch Route**:
    *   Update Load Balancer Rule:
        *   `if host == shop_X -> Proxy to Cluster B`
        *   `if host == shop_Y -> Proxy to Cluster B`
        *   `else -> Proxy to Cluster A`

### Step 3: Completion
*   Repeat until all shops are on Cluster B.
*   Decommission Cluster A.

---

## 3. Advantages (Solo Maintainer)
*   **Zero Code Logic**: Your Go code does NOT need to handle "If column exists". The code on Cluster B *knows* the column exists. The code on Cluster A *doesn't look for it*.
*   **Risk Isolation**: If `v1.1` is broken, only the migrated batch is affected. You can rollback (downgrade schema + switch route back).
*   **Zero Global Downtime**: The platform never goes down. Only individual shops might experience a blip during the switch.

## 4. Requirements
*   **Central DB**: Must be accessible by both clusters.
*   **Programmable LB**: Nginx, Caddy, or a simple Go Reverse Proxy that can look up "Routing Table" from Redis.
