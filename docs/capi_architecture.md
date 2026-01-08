# CAPI & Event Tracking Architecture (Hybrid)

## 1. Executive Summary
This document defines the architecture for the **Conversions API (CAPI)** integration for Meta (Facebook) and TikTok. The goal is to balance **Ad Optimization (Real-time)** with **Attribution Accuracy (Data Quality)**.

**Strategy:** "Local Brain, Central Muscle".
*   **Browser:** Sends events immediately (Real-time signal).
*   **Client Server:** Buffers events locally for session duration to enrich data.
*   **Central Server:** Receives batched events and handles delivery retry logic via RabbitMQ.

---

## 2. The Architecture

### 2.1 Component Overview

| Component | Location | Role | Tech Stack |
| :--- | :--- | :--- | :--- |
| **Browser Pixel** | User Device | Real-time signals (`fbp`, `fbc`, IP). | Javascript (GTM) |
| **Local Buffer** | App Node (Stateless) | Generates `event_id`. Buffers events in **Tenant Redis**. | Go + Redis |
| **External Gateway** | Third-Party Service | We post the batch to `CAPI_GATEWAY_URL`. | External HTTP |

### 2.2 The Forwarding Workflow
1.  **Event**: User views page.
2.  **Buffer**: App pushes event to Redis List `queue:capi`.
3.  **Worker**:
    *   Pops event.
    *   POSTs to `os.Getenv("CAPI_GATEWAY_URL")`.
    *   (No local RabbitMQ/Logic needed).
Instead of a strict 24-72h delay (which hurts ad learning), we use a **Session-Based Window** (Max 4 hours).

1.  **Guest Lands:**
    *   Server generates a unique `event_id` (UUID).
    *   Server stores `PageView` event in **Local Redis** (`capi:session:{id}`).
    *   Browser sends `PageView` to Meta immediately with this `event_id`.
2.  **User Browses:**
    *   `ViewContent`, `AddToCart` are added to the Local Redis list.
    *   Browser sends them immediately.
3.  **Enrichment Event (Purchase/Login):**
    *   User provides Name/Email/Phone.
    *   Local Server **Backfills** this user data to ALL previous events in the Redis list.
    *   *Result:* The `PageView` from 10 mins ago now has an Email attached.
4.  **Flush Trigger:**
    *   **Immediate:** On Purchase.
    *   **Timeout:** If no activity for 30 mins (Session End).
    *   **Max:** Force flush at 4 hours.
5.  **Dispatch:**
    *   Client VPS sends the *enriched batch* to the Central Server.

### 2.3 Deduplication Logic
*   **Browser:** Sends Event A at T+0s.
*   **Server:** Sends Event A at T+30m (Enriched).
*   **Meta:** Sees two events with same `event_id`. It *merges* them: taking the Timestamp from Browser (Speed) and User Data from Server (Quality).

---

## 3. Central Dispatch Infrastructure

### 3.1 Queue Technology Selection: RabbitMQ vs Kafka
We selected **RabbitMQ**.

| Feature | RabbitMQ | Kafka | Why RabbitMQ? |
| :--- | :--- | :--- | :--- |
| **Pattern** | Smart Broker / Job Queue | Dumb Broker / Streaming | CAPI is a "Job" (Deliver this payload), not a stream. |
| **Retries** | **Excellent (DLX)** | Difficult (Manual) | We need precise "Retry in 5m, then 15m" logic for API errors. |
| **Routing** | Flexible (Exchanges) | Fixed (Partitions) | Easy to route `meta_queue` vs `tiktok_queue`. |
| **Ack/Nack** | Per Message | Per Offset | We need to retry *specific* failed events, not block the whole queue. |

### 3.2 Dispatch Workflow
1.  **Ingest API:** Central Go API receives Batch. Validates format.
2.  **Queueing:** Pushes jobs to RabbitMQ `capi_events` exchange.
    *   Routing Key: `meta.purchase`, `tiktok.view`.
3.  **Workers:** Go Consumers allow `prefetch=100`.
    *   Attempt 1: POST to Meta Graph API.
    *   **Success:** Ack.
    *   **Failure (Rate Limit/500):** Nack -> Send to `retry_queue` with TTL (5m).
    *   **Failure (400 Bad Request):** Ack (Discard) -> Log Error.

---

## 4. Implementation Details

### 4.1 Schema (Event Context)
All events must carry:
```go
type CapiEvent struct {
    EventID      string `json:"event_id"`      // Critical for Dedup
    EventName    string `json:"event_name"`    // Purchase, ViewContent
    Timestamp    int64  `json:"timestamp"`
    ActionSource string `json:"action_source"` // "website"
    UserData     struct {
        Email    string `json:"em,omitempty"`  // Hashed (SHA256)
        Phone    string `json:"ph,omitempty"`  // Hashed
        FBP      string `json:"fbp"`
        FBC      string `json:"fbc"`
        ClientIP string `json:"client_ip"`
    }
}
```

### 4.2 Security
1.  **Encryption:** Client VPS -> Central Server must be HTTPS (TLS 1.3).
2.  **Auth:** Shared Secret Key or mTLS between Client VPS and Central Server.
3.  **Hashing:** PII (Email/Phone) must be SHA256 hashed *before* leaving the Client VPS Buffer.

---

## 5. Next Steps
1.  Implement **Local Redis Buffer** in `internal/modules/analytics`.
2.  Deploy **Central RabbitMQ** Cluster.
3.  Build **Central CAPI Worker** (Go).
