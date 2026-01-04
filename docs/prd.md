# Project Requirements Document (PRD)

**Project Name:** BizBundl E-commerce
**Type:** Single-Tenant E-commerce Platform
**Stack:** Go (Fiber), PostgreSQL, Templ, HTMX, Docker.

## 1. Executive Summary
A high-performance, single-tenant e-commerce solution designed for speed and simplicity. It allows a business to host their own isolated store with full data control. Ideally suited for the "Digital Product" market initially, with infrastructure to scale to physical goods.

## 2. Core Features (MVP - Digital First)

### A. Storefront (Customer View)
1.  **Product Discovery:**
    *   **Digital Products Only** (No shipping calculation needed).
    *   Instant Delivery (Email/Download Link) after payment.
    *   Homepage, PLP (Listing), PDP (Detail), Search.
2.  **Shopping Experience:**
    *   Cart & Checkout (Instant).
    *   **No Cash On Delivery (COD)** for MVP.
    *   Guest & Registered Checkout.
3.  **Authentication:**
    *   Customer Registration/Login.
    *   My Account: "My Downloads" / "Order History".

### B. Admin Panel (Merchant View)
1.  **Catalog Management:**
    *   Products: Title, Price, Description, **File Upload (for delivery)**.
    *   **Future (Physical):** Robust variation system options.
2.  **Payments (MVP Scope):**
    *   **Provider:** UddoktaPay (Sole provider for MVP).
    *   Merchant Config: API Key & URL in `payment_gateways`.
    *   **Couriers:** Integration with Pathao/RedX (Future/Infrastructure Ready).
    *   **RBAC:** Staff accounts with granular permissions (e.g., "See Orders" but not "Settings").
3.  **Order Management:**
    *   Status: Paid -> Completed (Auto).
    *   Manual overrides for support.

### C. Marketing & Intelligence (Growth Engine)
1.  **A/B Testing (PDP Optimizer):**
    *   **Mechanism:** Create multiple Layout variants for a single Product Page.
    *   **Traffic Split:** Randomly assign users (50/50) to Variant A or B (Sticky Session via Cookie).
    *   **Metrics:** Track Views vs. Conversions (Add to Cart / Purchase) to declare a "Winner".
2.  **Server-Side Tracking (Unified CAPI):**
    *   **Module:** One event bus (`Tracker.Send("Purchase", data)`).
    *   **Integrations:** Meta CAPI, TikTok Events API, GTM Server.
    *   **Logic:** Asynchronously sends events to all enabled providers to prevent latency.

## 3. Technical Architecture

### A. Performance & Caching Strategy
*   **Edge:** Cloudflare for HTML/Assets.
*   **App Configs (Read-Through Cache):**
    *   **Critical:** Configs (Payment Keys, Pixel IDs) are read 100x per second.
    *   **Mechanism:** App checks `MemoryStore` first. Miss -> Read DB -> Populate Memory.
    *   **Sync:** Updates to DB immediately invalidate/update the MemoryStore.
*   **Middleware:** Fiber `etag` and `cache` for Origin protection.

### B. Data Model & Isolation
*   **Database:** One standalone PostgreSQL database per store.
*   **Security:** Full isolation. Admin has complete control.

### B. High-Performance Page Builder (Strict Mode)
*   **Philosophy:** "Configuration over Design". Admins choose pre-compiled components and map data sources.
*   **Architecture:** Compiles `.templ` components. No runtime parsing.
*   **Safety:** Configuration validated at save time (e.g., ensuring a selected category actually exists).
*   **Reference:** See [Builder Architecture](docs/builder_architecture.md) for full design.

### C. Build & Deployment Strategy ("Source-Level Composition")
To support both mass-market users (Basic) and premium users (Custom Design) efficiently.

#### 1. The Build Philosophy
*   **Local Build / Remote Pull:** We build Docker images locally (or in CI) and push them to the registry. The production server **only pulls** and runs. It never builds.
*   **Baked-in Templates:** Templates (`.templ`) are compiled into the binary at build time, not mounted at runtime. Performance is prioritized.

#### 2. Pipeline A: The "Standard" Image (90% Users)
*   **Source:** `git/saas-core`.
*   **Process:** Standard `docker build` of the core repo.
*   **Result:** `ghcr.io/bizbundl/core:latest`.
*   **Update:** All "Basic" tier containers pull this single image.

#### 3. Pipeline B: The "Custom" Image (Pro Users)
*   **Source:** Custom Client Repo (`git/client-x`) + `saas-core:builder` base.
*   **Process (Injection Build):**
    1.  Pull `saas-core` (Builder stage).
    2.  **Delete** default `/views`.
    3.  **Copy** Client's custom `/views` into the build context.
    4.  Run `templ generate` (Compiles custom HTML into Go).
    5.  Run `go build`.
*   **Result:** A unique, optimized binary for that specific client.

## 4. Implementation Guidelines
*   **Code Style:** Strict SOLID/Clean Code.
*   **No Comments:** Self-documenting/meaningful naming.
*   **Performance:**
    *   **Cloudflare:** Heavy caching of assets.
    *   **Middleware:** Fiber `etag` and `cache` for Origin protection.
*   **Tech Stack:** Go 1.24, Fiber v2, Templ, SQLC.
