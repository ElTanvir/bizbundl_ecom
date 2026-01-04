package middleware

import (
	"bizbundl/token"
	"bizbundl/util"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Auth Middleware (Super Efficient / Stateless)
// Uses Paseto/JWT to verify user without DB lookup.
func Auth(tokenMaker token.Maker, userDuration, guestDuration, refreshThreshold time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get Token
		tokenString := c.Cookies("session_token")
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if len(authHeader) > 7 && strings.EqualFold(authHeader[0:7], "Bearer ") {
				tokenString = authHeader[7:]
			}
		}

		var payload *token.Payload
		var err error

		// 2. Verify Token (if exists)
		if tokenString != "" {
			payload, err = tokenMaker.VerifyToken(tokenString)
		}

		// 3. Create Guest Identity if missing or invalid
		// Note: We might want to separate "Strict Auth" (401) vs "Guest Identity" (Auto-create)
		// For E-commerce, we usually want Guest Identity everywhere by default.
		if tokenString == "" || err != nil {
			// If verification failed, clear bad cookie
			if err != nil {
				c.ClearCookie("session_token")
			}

			// Generate NEW Guest Identity
			guestID := uuid.New() // Using uuid direct
			// Actually payload needs UUID string
			newToken, newPayload, err := tokenMaker.CreateToken(guestID.String(), "guest", guestDuration)
			if err != nil {
				return util.APIError(c, fiber.StatusInternalServerError, fiber.NewError(fiber.StatusInternalServerError, "Failed to create guest session"))
			}

			// Set Cookie
			cookie := new(fiber.Cookie)
			cookie.Name = "session_token"
			cookie.Value = newToken
			cookie.Expires = time.Now().Add(guestDuration) // Guest duration
			cookie.HTTPOnly = true
			cookie.Secure = true
			cookie.SameSite = "Strict"
			c.Cookie(cookie)

			tokenString = newToken
			payload = newPayload
		}

		// 4. Set Context
		c.Locals("user_id", payload.UserID)
		c.Locals("user_role", payload.Role)

		// Auto-Renewal (Sliding Window) for Valid Tokens
		if time.Since(payload.IssuedAt) > refreshThreshold {
			// Generate NEW token
			duration := userDuration
			if payload.Role == "guest" {
				duration = guestDuration
			}

			newToken, _, err := tokenMaker.CreateToken(payload.UserID, payload.Role, duration)
			if err == nil {
				// Set Cookie
				cookie := new(fiber.Cookie)
				cookie.Name = "session_token"
				cookie.Value = newToken
				cookie.Expires = time.Now().Add(duration)
				cookie.HTTPOnly = true
				cookie.Secure = true
				cookie.SameSite = "Strict"
				c.Cookie(cookie)
			}
		}

		return c.Next()
	}
}
