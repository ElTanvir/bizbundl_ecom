package middleware

import (
	"bizbundl/token"
	"bizbundl/util"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Auth Middleware (Super Efficient / Stateless)
// Uses Paseto/JWT to verify user without DB lookup.
func Auth(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get Token
		tokenString := c.Cookies("session_token")
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if len(authHeader) > 7 && strings.EqualFold(authHeader[0:7], "Bearer ") {
				tokenString = authHeader[7:]
			}
		}

		if tokenString == "" {
			return util.APIError(c, fiber.StatusUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "Missing authentication token"))
		}

		// 2. Verify Token (CPU only, No DB)
		payload, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			c.ClearCookie("session_token")
			return util.APIError(c, fiber.StatusUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token"))
		}

		// 3. Set Context
		c.Locals("user_id", payload.UserID)

		// Auto-Renewal (Sliding Window)
		// Check if token is "old" enough to need a refresh (e.g. > 30 mins old)
		// Threshold: service.RefreshThreshold (we can't import service here closely, let's hardcode or pass config?
		// For efficiency, hardcode 30m here or move constants to a shared 'auth/types' pkg.
		// Let's hardcode 30m for now to avoid cycle, or move constants to 'token' pkg?
		// Let's assume 30m.
		if time.Since(payload.IssuedAt) > 30*time.Minute {
			// Generate NEW token
			// We need duration. If Role == "guest" -> 2 years. Else -> 2 hours.
			duration := 2 * time.Hour
			// if payload.Role == "guest" { duration = 2 * 365 * 24 * time.Hour }

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
