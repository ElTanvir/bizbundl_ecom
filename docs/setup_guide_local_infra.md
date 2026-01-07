# Local Infrastructure Setup Guide

This guide explains how to set up Redis and Elasticsearch locally using Docker, matching the default configuration of the BizBundl application.

## 1. Prerequisites
*   Docker & Docker Compose installed.
*   Go 1.22+ installed.

## 2. Start Infrastructure
We have provided a dedicated compose file for development infrastructure.

```bash
docker-compose -f docker-compose.dev-infra.yml up -d
```

This will start:
*   **Redis**: `localhost:6379` (No password, DB 0)
*   **Elasticsearch**: `localhost:9200` (No security, single node)

## 3. Configuration & Running

The application uses default values in `internal/config/config.go` that match this Docker setup.

### A. Default Mode (Redis Only)
By default, `ELASTIC_URL` is empty, so the app will use Postgres Trigram search (Fallback mode).

1.  Start the server:
    ```bash
    go run ./cmd/server/main.go
    ```
2.  Logs will show:
    ```
    INF Connected to Redis addr=localhost:6379
    WRN Elasticsearch URL not provided. Search will degrade to Database fallback.
    ```

### B. Full Mode (Redis + Elastic)
To enable Elasticsearch, you need to provide the `ELASTIC_URL` environment variable.

1.  Start the server with Env Var:
    ```bash
    export ELASTIC_URL="http://localhost:9200"
    go run ./cmd/server/main.go
    ```
2.  Logs will show:
    ```
    INF Connected to Redis addr=localhost:6379
    INF Connected to Elasticsearch url=http://localhost:9200
    ```

## 4. Multi-Tenant Configuration
The application supports multi-tenancy via Redis DB separation or Key Prefixes.

To change the default tenant settings:
```bash
export REDIS_DB=1
export REDIS_PREFIX="tenant_2:"
go run ./cmd/server/main.go
```

## 5. Troubleshooting
*   **Redis Connection Refused**: Ensure docker container is running (`docker ps`).
*   **Elastic Connection Refused**: Elasticsearch takes a few seconds to start. Wait for `Green status` or `started` logs in docker.
