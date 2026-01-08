package catalog

import (
	"bizbundl/internal/storefront/catalog/handler"
	"bizbundl/internal/storefront/catalog/service"
	"bizbundl/internal/server"
)

// Init initializes the Catalog module
func Init(app *server.Server) *service.CatalogService {
	svc := service.NewCatalogService(app.GetDB())
	handler := handler.NewCatalogHandler(svc)

	api := app.GetRouter().Group("/api/v1")
	handler.RegisterRoutes(api)
	return svc
}
