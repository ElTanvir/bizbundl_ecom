package cart

import (
	"bizbundl/internal/modules/cart/handler"
	"bizbundl/internal/modules/cart/service"
	"bizbundl/internal/server"
)

// Init initializes the Cart module
func Init(app *server.Server) {
	svc := service.NewCartService(app.GetDB())
	handler := handler.NewCartHandler(svc)

	api := app.GetRouter().Group("/api/v1")
	handler.RegisterRoutes(api)
}
