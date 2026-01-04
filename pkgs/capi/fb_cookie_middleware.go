package capi

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FbcCookieMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fbclid := c.Query("fbclid")
		if fbclid != "" && c.Cookies("_fbc") == "" {
			cookieValue := fmt.Sprintf("fb.1.%d.%s", time.Now().UnixMilli(), fbclid)

			cookie := &fiber.Cookie{
				Name:    "_fbc",
				Value:   cookieValue,
				Expires: time.Now().Add(90 * 24 * time.Hour),
				Path:    "/",
			}

			c.Cookie(cookie)
		}
		return c.Next()
	}
}
