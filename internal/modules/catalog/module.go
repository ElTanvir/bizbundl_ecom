package catalog

import (
	"bizbundl/internal/modules/catalog/handler"
	"bizbundl/internal/modules/catalog/service"
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
