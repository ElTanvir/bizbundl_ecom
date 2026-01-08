package middleware

import (
	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/infra/redis"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var validTenantID = regexp.MustCompile(`^[a-z0-9_]+$`)

// TenancyMiddleware wraps the request in a transaction with the correct search_path
func TenancyMiddleware(store db.DBStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Identify Tenant
		host := c.Hostname()
		// Simple subdomain extraction for MVP: shop1.localhost -> shop1
		// In production, use a proper mapping or Redis lookup
		tenantID := extractTenantID(host)

		// 2. Validate TenantID (Prevent SQL Injection)
		if !validTenantID.MatchString(tenantID) {
			// Fallback to public or error?
			// For now, default to 'public' for system routes, or error
			tenantID = "public"
		}

		// 3. Begin Transaction
		pool := store.GetPool()
		ctx := c.UserContext()
		tx, err := pool.Begin(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database Error")
		}

		// 4. Set Schema
		// Note: SET search_path cannot be easily parameterized in simple SQL.
		// Since we validated tenantID against strict regex, Sprintf is safe-ish.
		// Ideally use QuoteIdentifier.
		query := fmt.Sprintf("SET search_path TO %s, public", tenantID)
		if _, err := tx.Exec(ctx, query); err != nil {
			tx.Rollback(ctx)
			return c.Status(fiber.StatusInternalServerError).SendString("Schema Error")
		}

		// 5. Inject Tx into Context
		ctxWithTx := context.WithValue(ctx, db.TxKey, tx)
		ctxWithTenant := context.WithValue(ctxWithTx, redis.TenantKey, tenantID)
		c.SetUserContext(ctxWithTenant)

		// 6. Inject TenantID for Redis keys (Fiber Locals for Handler access)
		c.Locals("tenant_id", tenantID)

		// 7. Next Handler
		if err := c.Next(); err != nil {
			tx.Rollback(ctx)
			return err
		}

		// 8. Commit
		// Only commit if status code is success (2xx/3xx)
		if c.Response().StatusCode() >= 400 {
			tx.Rollback(ctx)
		} else {
			if err := tx.Commit(ctx); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Commit Error")
			}
		}

		return nil
	}
}

func extractTenantID(host string) string {
	// Strip port if present
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	// Handle Development / Localhost defaults
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return "public"
	}

	parts := strings.Split(host, ".")
	if len(parts) > 0 {
		// If it's an IP address (numeric), likely we should treat as public or handle differently
		// For MVP, assuming subdomains are alphabetic.
		// "192.168.1.1" -> "192" -> valid regex? yes.
		// Let's explicitly check if it looks like a shop subdomain.
		// For now, the explicit check above handles local dev.
		return parts[0]
	}
	return "public"
}
