# CAPI Implementation Plan (Phase 5)

## Goal
Decouple CAPI event processing from the Request/Response cycle to improve performance and reliability.

## Architecture
1.  **Ingestion**: `SendPageViewEvent` (and others) will **RPUSH** the event payload to a Redis List (`queue:capi:{tenant_id}` or global `queue:capi` with tenant metadata).
2.  **Processing**: A background worker (Goroutine or dedicated process) will **BLPOP** from the list.
3.  **Dispatch**: The worker sends the HTTP request to Meta/TikTok.

## Changes

### 1. Redis Queue Helper (`internal/infra/redis`)
- Add `Enqueue(ctx, key, value)`
- Add `Dequeue(ctx, key)`

### 2. CAPI Refactor (`pkgs/capi`)
- Modify `sendPageViewEventAsync` to `EnqueueCAPIEvent`.
- Serialize `CAPIEvent` to JSON.
- Push to Redis.

### 3. Analytics Worker (`cmd/analytics_worker`)
- New CLI tool.
- Infinite Loop:
    - `Redis.BLPOP`
    - Unmarshal Event.
    - `sendCAPIRequest`
