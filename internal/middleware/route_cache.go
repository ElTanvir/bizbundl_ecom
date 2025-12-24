package middleware

import (
	"portfolioed/internal/store"
	"portfolioed/util"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func RouteCacheMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if util.IsCacheableRoute(c) {
			return c.Next()
		}
		entry, found := store.GetRouteBytes(c.Path())
		if !found {
			return c.Next()
		}
		log.Info().Str("Etag", entry.ETag).Str("path", c.Path()).Msg("Serving from route cache")
		return util.ServeWithETag(c, entry.Data, entry.ETag)
	}
}
