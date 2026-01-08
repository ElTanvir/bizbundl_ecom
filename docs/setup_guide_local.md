# Local Development Setup

## 1. Prerequisites
*   Go 1.22+
*   Docker & Docker Compose

## 2. Start Infrastructure
```bash
docker-compose -f docker-compose.dev-infra.yml up -d
```
Starts Redis (6379) and Postgres (default).

## 3. Running the App
The app defaults to `localhost` config.

```bash
go run ./cmd/server/main.go
```

## 4. Simulating Multi-Tenancy Locally
To test Schema Isolation:

1.  **Hosts File**: Map domains to localhost.
    *   `127.0.0.1 shop1.local`
    *   `127.0.0.1 shop2.local`
2.  **Seeding**:
    *   Run the seeder to create `shop_1` and `shop_2` schemas.
    *   (Seeder tool TBD in Refactoring phase).
3.  **Access**:
    *   Visit `http://shop1.local:8080` -> Hits Schema 1.
    *   Visit `http://shop2.local:8080` -> Hits Schema 2.
