package frontend

import (
	cartservice "bizbundl/internal/modules/cart/service"
	"bizbundl/internal/modules/catalog/service"
	"bizbundl/internal/modules/page_builder"
	"bizbundl/internal/server"
	"bizbundl/internal/views/frontend/handler"
)

func Init(app *server.Server) {
	catalogSvc := service.NewCatalogService(app.GetDB())
	cartSvc := cartservice.NewCartService(app.GetDB())
	pbModule := page_builder.Init(app, catalogSvc, cartSvc)
	h := handler.NewFrontendHandler(catalogSvc, cartSvc, pbModule.Service, pbModule.Resolver)

	// Frontend Routes
	// Serve static assets if needed, but usually handled by Fiber static or Nginx
	// app.GetRouter().Static("/static", "./public")

	// HTML Pages
	routes := app.GetRouter().Group("/")
	routes.Get("/", h.HomePage)
	routes.Get("/product/:slug", h.ProductPage)
	// Dynamic Landing Pages (Catch-All) - Must be last!
	routes.Get("/*", h.RenderLandingPage)

	// Cart
	routes.Get("/cart", h.CartPage)
	routes.Post("/cart/add", h.AddToCart) // Simplified non-REST for HTMX ease or keep REST? HTMX usually POST.
	routes.Post("/cart/update", h.UpdateCartItem)
	routes.Delete("/cart/items/:id", h.RemoveCartItem)
}
