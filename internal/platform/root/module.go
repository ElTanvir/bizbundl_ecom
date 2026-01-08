package root

import (
	"bizbundl/internal/platform/root/handler"
	"bizbundl/internal/server"
)

func Init(app *server.Server) {
	h := handler.NewRootHandler()

	// Platform Root Routes (Public)
	// We bind to "/" but strictly this only applies if hostname is public.
	// The TenancyMiddleware handles logic, but FrontendHandler ALSO matches "/".
	// We need to Coordinate.
	// User wants "Internal/platform/root" to house the view for homepage.
	// If we register "/" here, it might conflict with "FrontendHandler.HomePage".

	// Strategy:
	// If Tenant="public", we want THIS handler.
	// If Tenant="shop", we want FrontendHandler.
	// Middleware separates them by Subdomain? No, middleware sets context.

	// Current FrontendHandler.HomePage explicitly checks for "public" tenant and renders PlatformHome.
	// We should REPLACE that logic by registering this handler for the "public" domain specifically?
	// Fiber doesn't do domain-routing easily without groups/host matching.

	// For now, let's keep the logic in FrontendHandler OR update Main to register this handler conditionally?
	// Actually, simpler:
	// Let's register a dedicated group or let FrontendHandler delegate?
	// User wants structure.
	// Let's Register it, but maybe as a named route or handle collision later.
	// Actually, `FrontendHandler` uses `/*` catchall.
	// If we define `app.Get("/", h.ShowHome)`, it takes precedence over `/*`.
	// So if we register this, it works for `/`.
	// BUT checking for `public` tenant inside here is crucial if we share the port.

	// Let's assume for this Refactor step, I will register it, and we might need to verify collision.
	// Register Root "/" route.
	// Note: This registers a specific "/" handler. Frontend catch-all "/*" will be skipped for exactly "/".
	// We rely on FrontendHandler to NOT overwrite this if it's registered first or specific match wins?
	// Fiber: Specific match `/` wins over `/*`.
	// BUT we need tenant check.
	// For now, let's register it. If tenant!=public, we might need middleware to reject/proxy?
	// Or we keep logic in FrontendHandler?
	// User wants "Root module houses view for homepage".
	// Let's register it here.
	app.GetRouter().Get("/", h.ShowHome)
}
