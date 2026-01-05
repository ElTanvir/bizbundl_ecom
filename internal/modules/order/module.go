package order

import (
	cartService "bizbundl/internal/modules/cart/service"
	catalogService "bizbundl/internal/modules/catalog/service"
	"bizbundl/internal/modules/order/handler"
	"bizbundl/internal/modules/order/service"
	"bizbundl/internal/modules/payment/providers/uddoktapay"
	"bizbundl/internal/server"
)

type Module struct {
	handler *handler.OrderHandler
	service *service.OrderService
}

func Init(app *server.Server, cartSvc *cartService.CartService, catalogSvc *catalogService.CatalogService) *Module {
	svc := service.NewOrderService(app.GetDB())

	// Payment GW
	pgw := uddoktapay.New("") // Uses default Sandbox Key internaly

	h := handler.NewOrderHandler(cartSvc, svc, catalogSvc, pgw)

	// Register Routes
	g := app.GetRouter().Group("/order")
	g.Post("/checkout", h.Checkout)
	g.Get("/checkout", h.ShowCheckoutPage) // New Route
	g.Get("/payment/callback", h.PaymentCallback)
	g.Get("/success/:id", h.SuccessPage)

	return &Module{
		handler: h,
		service: svc,
	}
}
