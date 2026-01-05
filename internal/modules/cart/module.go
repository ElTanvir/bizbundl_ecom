package cart

import (
	"bizbundl/internal/modules/cart/handler"
	"bizbundl/internal/modules/cart/service"
	"bizbundl/internal/server"
)

// Init initializes the Cart module
func Init(app *server.Server) *service.CartService {
	svc := service.NewCartService(app.GetDB())
	handler := handler.NewCartHandler(svc)

	api := app.GetRouter().Group("/api/v1")
	handler.RegisterRoutes(api)

	// Frontend Actions (HTMX/Form)
	// We register specific actions at root level to support frontend forms
	app.GetRouter().Post("/cart/items", handler.AddItem)

	return svc
}
