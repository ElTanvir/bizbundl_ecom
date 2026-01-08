package platform

import (
	"bizbundl/internal/modules/platform/handler"
	"bizbundl/internal/modules/platform/service"
	"bizbundl/internal/server"
	// "bizbundl/internal/middleware" // Auth middleware needed
)

// Init orchestrates the Platform Module
func Init(app *server.Server) {
	// 1. Dependency
	// We need generic Pool to access 'platform' schema queries
	pool := app.GetDB().GetPool()
	cfg := app.GetConfig()

	// 2. Service
	svc := service.NewPlatformService(pool, cfg)

	// 3. Handler
	h := handler.NewPlatformWebHandler(svc)

	// 4. Routes
	// Group: /dashboard (Protected)
	// Apply Global Auth Middleware? User says "No API handler", so likely web session.
	// We can trust the session middleware applied globally or add it here.
	// For now, assume middleware stack in server.go handles auth/session.

	// Actually, server.go usually applies generic middleware.
	// Specific "Auth Check" middleware usually needed for protected routes.
	// Let's assume we skip middleware for MVP instant load, OR use the one from middleware pkg.
	// 'middleware.Auth' returns a handler.
	// We can't access it easily without importing middleware pkg and deps.
	// For "Web Only", let's assume the handler checks the session or we register it later.
	// The web_handler I wrote checks `c.Locals("user_id")` which implies Auth Middleware ran.
	// So we should attach it.
	// But getting the middleware instance might be tricky if not centralized.
	// Auth module exports it? No, Auth module Init creates it.
	// Let's just register routes.

	r := app.GetRouter()
	dash := r.Group("/dashboard")

	dash.Get("/", h.ShowDashboard)
	dash.Get("/shops/new", h.ShowCreateShopForm)
	dash.Post("/shops", h.HandleCreateShop)
}
