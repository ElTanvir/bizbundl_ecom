package auth

import (
	"bizbundl/internal/constants"
	"bizbundl/internal/middleware"

	// "bizbundl/internal/middleware" // Global?
	"bizbundl/internal/platform/auth/handler"
	"bizbundl/internal/platform/auth/service"

	// cartservice "bizbundl/internal/storefront/cart/service"
	"bizbundl/internal/server"
)

// Init initializes the Auth module: wires the Service/Handler and registers routes.
func Init(app *server.Server) {
	svc := NewAuthService(app)
	h := handler.NewAuthHandler(svc)

	// Auth Middleware (Global)
	// We pass the token maker directly to middleware
	// We also pass session durations configuration
	authMiddleware := middleware.Auth(
		app.GetTokenMaker(),
		constants.UserSessionDuration,
		constants.GuestSessionDuration,
		constants.RefreshThreshold,
	)

	// Apply Auth Middleware Globally (to ensure Guest Session on all routes)
	app.GetRouter().Use(authMiddleware)

	// View Routes (HTMX Pages & Actions)
	app.GetRouter().Get("/login", h.ShowLoginForm)
	app.GetRouter().Post("/login", h.Login)

	// Logout Action
	app.GetRouter().Post("/logout", h.Logout)

	// Register Page (TODO: Create Register View)
	app.GetRouter().Post("/register", h.Register)

	// Protected
	// api.Get("/me", authMiddleware, h.Me) // 'Me' is likely used by UI to get state? Or we use template data.
	// For now, let's keep 'Me' disabled or move to web route if needed.

}

func NewAuthService(app *server.Server) *service.AuthService {
	// We need 'platform' queries.
	// app.GetDB() returns Tenant Store.
	// We need to construct Platform Queries manually or via helper.
	pool := app.GetDB().GetPool()
	queries := service.NewPlatformQueries(pool) // We'll add this helper or inline it in service pkg
	return service.NewAuthService(queries, app.GetTokenMaker())
}
