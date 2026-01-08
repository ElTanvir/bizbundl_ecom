# Capacity Planning: Schema-Based Architecture

**Hardware Reference**: 6 vCPU / 12 GB RAM / NVMe SSD ($12/mo VPS).

## 1. Resource Efficiency
With Schema-Based Multi-Tenancy, resource usage is extremely efficient because tenants share the same memory space (Postgres Shared Buffers) and Connection Pool.

| Resource | Usage Pattern | Limit |
| :--- | :--- | :--- |
| **Idle Memory** | ~0 MB per tenant (Shared) | **Limited by Disk Space** |
| **Active Memory** | Depends on concurrent reqs | ~10k concurrent reqs |
| **Connections** | Shared Pool (50-100 conns total) | DB Max Connections |

## 2. Capacity Estimates (Single Application Node)

### A. Tenant Count (Storage Bound)
*   If average shop usage = 100MB (Images/Data).
*   Server Disk = 100GB.
*   **Capacity = ~1,000 Shops.** (Limited by Disk).
*   *Solution: Offload uploads to S3/Object Storage to unlock infinite density.*

### B. Tenant Count (Traffic Bound)
*   Throughput: Go+Fiber is extremely fast (>20k req/sec hello world, ~2k req/sec real app).
*   **Scenario: "Sleepy SaaS"** (90% low traffic, 10% active):
    *   **Can serve ~5,000+ Tenants comfortably.**

## 3. Scaling Strategy (Database Splitting)

When the "Mother Node" (DB) hits CPU/IOPS limits (approx 10k tenants), we split the database layer.

**Mechanism: Registry-Based Routing**
1.  **Architecture**:
    *   **Registry**: A lookup table (Redis/Config) mapping `TenantID -> DatabaseConnectionURL`.
    *   **DB 1 (Mother)**: Hosts Tenants 1-10,000.
    *   **DB 2 (Child)**: Hosts Tenants 10,001+.
2.  **Workflow**:
    *   Tenant `shop_new` registers.
    *   System checks capacity. "DB 1 Full".
    *   System provisions Schema on **DB 2**.
    *   Registry updated: `shop_new -> DB_2_URL`.
3.  **Code Impact**:
    *   Middleware looks up *Connection String* before `SET search_path`.
    *   Instead of 1 Global Pool, the App maintains a Map of Pools (`DB_1 -> Pool`, `DB_2 -> Pool`).
    *   *Note: This is "Application-Side Sharding".*

**Legacy Handling**: No need to migrate data. Old tenants stay on DB 1 forever until you choose to move them.

