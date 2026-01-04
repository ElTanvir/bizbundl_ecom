package catalog

import (
	"bizbundl/internal/server"
)

// Init initializes the Catalog module
func Init(app *server.Server) {
	svc := NewCatalogService(app.GetDB())
	handler := NewCatalogHandler(svc)

	api := app.GetRouter().Group("/api/v1")
	handler.RegisterRoutes(api)
}
