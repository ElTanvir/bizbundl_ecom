package util

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"portfolioed/internal/store"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	err := component.Render(c.Context(), buf)
	if err != nil {
		return err
	}
	if !IsCacheableRoute(c) {
		return c.Send(buf.Bytes())
	}
	renderedHTML := make([]byte, buf.Len())
	copy(renderedHTML, buf.Bytes())
	etag := generateETag(renderedHTML)
	go store.SetRouteBytes(c.Path(), renderedHTML, etag)
	// capi.SendPageViewEvent(c)
	return ServeWithETag(c, renderedHTML, etag)
}

func ServeWithETag(c *fiber.Ctx, html []byte, etag string) error {
	c.Set("ETag", etag)
	c.Set("Cache-Control", "public, max-age=3600, must-revalidate")
	c.Set("CDN-Cache-Control", "public, max-age=86400, stale-while-revalidate=604800")
	c.Set("Vary", "HX-Request, HX-Boosted")
	if match := c.Get("If-None-Match"); match == etag {
		c.Set("X-Cache-Status", "NOT-MODIFIED")
		return c.SendStatus(fiber.StatusNotModified)
	}
	c.Set("X-Cache-Status", "HIT")
	return c.Send(html)
}

func generateETag(data []byte) string {
	h := fnv.New64a()
	h.Write(data)
	return fmt.Sprintf(`"%x"`, h.Sum64())
}

func IsCacheableRoute(c *fiber.Ctx) bool {
	if c.Method() != fiber.MethodGet {
		return false
	}
	path := c.Path()
	// Admin routes protection: /admin, /admin/*, /adm, /backend, /api/admin, etc.
	return !isAdminLikeRoute(path)
}

func isAdminLikeRoute(path string) bool {
	return strings.Contains(path, "/admin")
}
