package auth

import (
	"bizbundl/internal/constants"
	"bizbundl/internal/middleware"
	"bizbundl/internal/storefront/auth/handler"
	"bizbundl/internal/storefront/auth/service"
	cartservice "bizbundl/internal/storefront/cart/service"
	"bizbundl/internal/server"
)

// Init initializes the Auth module: wires the Service/Handler and registers routes.
func Init(app *server.Server) {
	svc := NewAuthService(app)
	// Initialize CartService for linking (Ideally, this should be a singleton in app, but new instance is fine)
	cartSvc := cartservice.NewCartService(app.GetDB())
	h := handler.NewAuthHandler(svc, cartSvc)

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

	// Routes
	api := app.GetRouter().Group("/api/v1/auth")

	// Public
	api.Post("/register", h.Register)
	api.Post("/login", h.Login)
	api.Post("/logout", h.Logout)

	// Protected
	api.Get("/me", authMiddleware, h.Me)
}

func NewAuthService(app *server.Server) *service.AuthService {
	return service.NewAuthService(app.GetDB(), app.GetTokenMaker())
}
