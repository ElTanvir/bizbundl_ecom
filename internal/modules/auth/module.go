package auth

import (
	"bizbundl/internal/server"
)

// Init initializes the Auth module: wires the Service/Handler and registers routes.
func Init(app *server.Server) {
	svc := NewAuthService(app.GetDB())
	handler := NewAuthHandler(svc)

	// Auth routes are usually public except "me" which needs middleware (later)
	api := app.GetRouter().Group("/api/v1")
	handler.RegisterRoutes(api)
}
