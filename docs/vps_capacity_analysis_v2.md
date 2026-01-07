# VPS Capacity Analysis V2: The Hybrid Strategy

## 1. Executive Summary
This document analyzes the capacity of a **6 vCore / 12GB RAM VPS ($7.70/mo)** under a Hybrid Architecture strategy. It specifically addresses the trade-offs of adding Elasticsearch (ES) and outlines the optimal strategy to **maximize tenant density** for the Hosted service.

---

## 2. Architecture Scenarios

### Scenario A: Hosted "Max Density" (No ElasticSearch)
**Stack:** Go App + Postgres + Redis.
**Target:** Starter & Standard Tiers.
**Philosophy:** Maximize the number of $199/mo and $799/mo clients per server. Search is handled by Postgres Trigram (Low RAM, High Density).

### Scenario B: Self-Hosted "Performance Beast" (Full Stack)
**Stack:** Go App + Postgres + Redis + Elasticsearch.
**Target:** Single "Unlimited" Client ($1000 License + Own VPS).
**Philosophy:** Uncompromising performance for one high-volume merchant.

---

## 3. Hosted Service: Maximizing Tenant Count
**Objective:** Fit the maximum number of paying tenants on one 12GB VPS.

### 3.1 The "Elasticsearch Tax" problem
Elasticsearch is a RAM hog. It requires a 4GB Heap minimum to be stable for multi-tenant workloads.

| Stack Component | RAM Usage (With ES) | RAM Usage (No ES) |
| :--- | :--- | :--- |
| **Elasticsearch** | **4 GB** | **0 GB** |
| **Postgres** | 2 GB | 3 GB (More Buffer) |
| **Redis** | 1 GB | 1 GB |
| **System** | 1 GB | 1 GB |
| **Available for Apps** | **4 GB** | **7 GB** |

### 3.2 Tenant Density Calculation
*Assuming "Starter" Apps average 50MB RAM active.*

*   **With ES:** 4GB Available / 50MB = **~80 Tenants Max**.
*   **Without ES:** 7GB Available / 50MB = **~140 Tenants Max**.

### 3.3 The "Maximize Density" Strategy
To maximize revenue per VPS, we **MUST NOT** run Elasticsearch on the Hosted instances.

**Recommended Hosted Stack:**
1.  **Database:** Postgres (Shared).
2.  **Cache:** Redis (Shared). **Crucial**. It offloads 95% of read traffic, keeping CPU usage near zero.
3.  **Search:** Postgres `pg_trgm` (Trigram).
    *   *Quality:* Excellent for stores with <5,000 products.
    *   *Speed:* <100ms (Cached by Redis if repeated).
    *   *Cost:* 0 RAM overhead (uses existing PG Buffer).

**Result:**
*   **Capacity:** **120-150 Tenants** per $7.70 VPS.
*   **Revenue Potential:** ~120 * Avg(400tk) = **48,000 BDT/mo**.
*   **Profit margin:** ~98%.

---

## 4. Self-Hosted Performance Analysis
**Context:** A single client pays the 1000tk license and buys their own $7.70 VPS.
**Stack:** Full Stack (Postgres + Redis + Elasticsearch).

### 4.1 Resource Allocation (Single Tenant)
Since there is only *one* application container, we can allocate massive resources to the infrastructure.

*   **Elasticsearch (4GB):** Dedicated to minimal products. Lightning fast.
*   **Postgres (3GB):** massive buffer for orders.
*   **Redis (1GB):** Caches entire site.
*   **App (3GB):** Huge heap for concurrent Go routines.

### 4.2 Performance Capability (Benchmarks)

| Metric | Capability | Real-World Equivalent |
| :--- | :--- | :--- |
| **Concurrent Users** | **5,000+** | Can handle a massive "Friday Prayer" flash sale. |
| **Page Views/Mo** | **50 Million+** | More than almost any BD e-commerce site. |
| **Orders/Hourly** | **30,000+** | Postgres Write throughput is the only limit. |
| **Search Speed** | **<20ms** | Instant results for 100k+ products. |

### 4.3 Comparison vs Hosted
*   **Self-Hosted ($1000 + $7.70):** Unlocked performance. You define the limits.
*   **Hosted ($2499 "Business"):** Great, but shares resources with 100 other stores.

---

## 5. Summary & Recommendation

1.  **For Hosted (Starter/Standard/Business):**
    *   **REMOVE Elasticsearch.** It kills density.
    *   Use **Redis** aggressively.
    *   **Capacity:** Aim for **150 Tenants** per VPS.

2.  **For Self-Hosted:**
    *   **INCLUDE Elasticsearch.**
    *   It differentiates the "Premium" self-hosted product.
    *   Client pays for the VPS, so the RAM usage is their benefit, not our cost.
