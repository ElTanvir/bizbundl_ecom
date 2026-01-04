package cart

import (
	"bizbundl/internal/server"
)

// Init initializes the Cart module
func Init(app *server.Server) {
	svc := NewCartService(app.GetDB())
	handler := NewCartHandler(svc)

	api := app.GetRouter().Group("/api/v1")
	handler.RegisterRoutes(api)
}
