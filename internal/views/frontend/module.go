package frontend

import (
	"bizbundl/internal/modules/catalog/service"
	"bizbundl/internal/server"
	"bizbundl/internal/views/frontend/handler"
)

func Init(app *server.Server) {
	catalogSvc := service.NewCatalogService(app.GetDB())
	h := handler.NewFrontendHandler(catalogSvc)

	// Frontend Routes
	// Serve static assets if needed, but usually handled by Fiber static or Nginx
	// app.GetRouter().Static("/static", "./public")

	// HTML Pages
	routes := app.GetRouter().Group("/")
	routes.Get("/", h.HomePage)
}
