# VPS Capacity & Pricing Strategy Analysis

## 1. Hardware Analysis: The "Value King" ($7.70/mo)
We are analyzing the **6 vCore / 12GB RAM** server.

| Spec | Value | Notes |
| :--- | :--- | :--- |
| **Price** | **$7.70 / mo** (~920 BDT) | Extremely low entry barrier. |
| **RAM** | **12 GB** | The primary constraint for tenant density. |
| **Capacity** | **~40 Tenants** | Safe allocation of 200MB/tenant (Real usage <50MB). |
| **Cost/Tenant** | **~23 BDT** | Hardware cost is negligible. |

---

## 2. New 4-Tier Pricing Strategy
**Logic:** Solve the "Under-priced Growth" issue by inserting a middle tier and raising the top tier.

| Feature | **Starter (Micro)** | **Standard (SME)** | **Business (Pro)** | **Self-Hosted** |
| :--- | :--- | :--- | :--- | :--- |
| **Target** | New Store | Growing Store | Established Brand | Technical / High Vol |
| **Price** | **৳199 / mo** | **৳799 / mo** | **৳2,499 / mo** | **৳1,000 / mo** |
| **Orders/Mo** | **300** | **1,500** | **5,000** | **Unlimited** |
| **Products** | 100 | 1,000 | 5,000 | **Unlimited** |
| **Storage** | 500 MB | 2 GB | 10 GB | **Own Server** |
| **Value** | *Entry Drug* | *The "Sweet Spot"* | *High Margin* | *License Fee Only* |

### Why this structure?
1.  **Starter (199tk):** Still the "No Brainer" vs Competitor (500tk). Captures volume.
2.  **Standard (799tk):** The **New Middle Tier**.
    *   Target: 1,500 Orders.
    *   Price: Affordable, but healthy profit.
    *   Competitor charges ~1500tk+ for this volume.
3.  **Business (2,499tk):** The **Adjusted Growth Tier**.
    *   Target: 5,000 Orders (High volume).
    *   We stopped under-pricing this. 2,499tk is fair for a business doing 5k orders/mo.
4.  **Self-Hosted (1,000tk):**
    *   **The "Unlimited" Solution.**
    *   If a client complains about the 5,000 limit on Business tier, we upsell them to Self-Hosted.
    *   "You want unlimited? Pay us 1,000tk license and buy your own $15 VPS."
    *   **Benefit:** We take 0 hardware risk. They get unlimited power.

---

## 3. Financial Projections (Per 40-User Server)

**Scenario: A Healthy Mix**
*   **20 Starters** (@ 199): 3,980 BDT
*   **15 Standards** (@ 799): 11,985 BDT
*   **5 Business** (@ 2,499): 12,495 BDT
*   **Total Revenue:** **~28,460 BDT / month**
*   **Server Cost:** ~920 BDT
*   **Net Profit:** **~27,540 BDT / month** per server.

### Profit Multipliers
*   **1 Server (40 Clients):** 27.5k Net Profit.
*   **5 Servers (200 Clients):** ~1.3 Lakh Net Profit.
*   **+ 50 Self-Hosted Licenses:** +50k Pure Profit.

## 4. Operational Strategy
1.  **Strict Limits:** The `Starter` tier MUST have hard limits (300 orders) enforced by code. This forces the upgrade to `Standard`.
2.  **Upsell Path:**
    *   User hits 300 orders -> Upgrade to Standard (799tk).
    *   User hits 1,500 orders -> Upgrade to Business (2499tk).
    *   User wants "Unlimited" -> "Switch to Self-Hosted, it's cheaper (1000tk license) + you own the server."
