package order

import (
	cartService "bizbundl/internal/modules/cart/service"
	"bizbundl/internal/modules/order/handler"
	"bizbundl/internal/modules/order/service"
	"bizbundl/internal/server"
)

type Module struct {
	handler *handler.OrderHandler
	service *service.OrderService
}

func Init(app *server.Server, cartSvc *cartService.CartService) *Module {
	svc := service.NewOrderService(app.Store)
	h := handler.NewOrderHandler(cartSvc, svc)

	// Register Routes
	g := app.App.Group("/order")
	g.Post("/checkout", h.Checkout)
	g.Get("/success/:id", h.SuccessPage)

	return &Module{
		handler: h,
		service: svc,
	}
}
