package auth

import (
	"bizbundl/internal/constants"
	"bizbundl/internal/middleware"
	"bizbundl/internal/modules/auth/handler"
	"bizbundl/internal/modules/auth/service"
	"bizbundl/internal/server"
)

// Init initializes the Auth module: wires the Service/Handler and registers routes.
func Init(app *server.Server) {
	svc := service.NewAuthService(app.GetDB(), app.GetTokenMaker())
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

	// Routes
	api := app.GetRouter().Group("/api/v1/auth")

	// Public
	api.Post("/register", h.Register)
	api.Post("/login", h.Login)
	api.Post("/logout", h.Logout)

	// Protected
	api.Get("/me", authMiddleware, h.Me)
}
