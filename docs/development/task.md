# Task Checklist - BizBundl E-Commerce

## Phase 0: Infrastructure & Refactoring (Priority)
- [/] **Analysis**: Dependency and Config analysis.
- [x] **Config**: Add Redis & Elastic support to `internal/config`.
- [x] **Infra**: Implement Redis Client (Mandatory).
- [x] **Infra**: Implement Elasticsearch Client (Optional/Fallback).
- [x] **Refactor**: Move `page_builder` to `pkgs/`.
- [x] **Refactor**: Replace `go-cache` with Redis for Page Caching.
- [x] **Docs**: Create Local Infra Setup Guide (Docker + Config).
- [x] **Docs**: Establish ADR Framework (`docs/adr`).
- [ ] **Testing**: Achieve High Test Coverage (Iterative).

## Phase 1: Architecture Refactor (Schema Multi-Tenancy)
- [x] **Config**: Update `config.go` (Remove Multi-DB, Add Schema Settings).
- [x] **Middleware**: Implement `TenancyMiddleware` (Host Parsing).
- [x] **Database**: Implement `SchemaInjector` (SET search_path).
- [x] **Redis**: Implement `NamespaceWrapper` (Prefix Keys).
- [x] **Migration**: Implement `cmd/migrate_worker` (Multi-Schema Iteration).

## Phase 2: Core Features (Existing)
- [x] Initial Repo Setup (Go, Fiber, Templ, Docker)
- [x] Database Schema & SQLC Setup
- [x] Auth Module (Login/Register/Sessions)

## Phase 2: Catalog & Content
- [x] Product Module (CRUD, SEO Slugs)
- [x] Category Module
- [x] Page Builder Core (Registry, Components)
- [x] Component: Product Grid (Variants)
- [x] Component: Hero/Banner
- [/] **Documentation**: VPS Capacity Analysis

## Phase 3: Cart & Order (Basic)
- [x] Cart Module (Guest/User Merging)
- [x] Order Module (Schema, Basic Creation)
- [ ] Order Management (Admin View)

## Phase 4: Page Builder & Routing
- [x] Dynamic Routing (catch-all `/*`)
- [x] Registry Duplicate Protection
- [ ] Checkout Widget Component

## Phase 5: Platform & Shop Management (New MVP Priority)
- [x] **Schema**: Add `public.shops` and `public.subscriptions` tables.
- [ ] **Feature**: Shop Creation Flow (Provisoning Schema).
- [ ] **Feature**: Owner Dashboard (List Shops).

## Phase 6: CAPI (External Integration)
- [ ] **Config**: Add `CAPI_GATEWAY_URL`.
- [ ] **Client**: Refactor `analytics_worker` to hit External Gateway.

## Phase 7: Payments (Billing)
- [ ] **Provider**: UddoktaPay / Stripe.
- [ ] **Flow**: Subscription Payments for Shops.
