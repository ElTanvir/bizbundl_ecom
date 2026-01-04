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
		c.Locals("user_role", payload.Role)

		// Auto-Renewal (Sliding Window)
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
