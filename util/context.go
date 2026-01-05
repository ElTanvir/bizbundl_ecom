package util

import (
	"bizbundl/internal/config"
	"bizbundl/token"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// GetUserIDFromContext retrieves the UserID from the Fiber context (JWT/Paseto)
func GetUserIDFromContext(c *fiber.Ctx) pgtype.UUID {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return pgtype.UUID{Valid: false}
	}
	var id pgtype.UUID
	id.Scan(payload.UserID) // Payload usually has String UUID
	return id
}

// GetSessionIDFromContext retrieves the SessionID from the Fiber context (Cookie)
func GetSessionIDFromContext(c *fiber.Ctx) pgtype.UUID {
	cookie := c.Cookies(config.SessionCookieName)
	if cookie == "" {
		return pgtype.UUID{Valid: false}
	}
	// We assume the cookie IS the session ID or a token resolving to it?
	// In strict architecture, we might look up "sessions" table by token.
	// Converting cookie string (token) to UUID directly is risky if it's not a UUID.
	// However, if we store Session ID (UUID) in cookie, this works.
	// If we store a Token, we need to lookup using a service.

	// For MVP, assuming "session_id" cookie holds the UUID.
	// If it holds "session_token", we can't just cast to UUID.
	// Let's assume for now it holds the UUID string.

	var id pgtype.UUID
	err := id.Scan(cookie)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return id
}
