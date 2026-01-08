package page_builder

import (
	cart_service "bizbundl/internal/storefront/cart/service"
	catalog_service "bizbundl/internal/storefront/catalog/service"
	"bizbundl/internal/server"
	"bizbundl/pkgs/page_builder/resolver"
	"bizbundl/pkgs/page_builder/service"

	// Import Atomic Components
	"bizbundl/pkgs/components/checkout"
	_ "bizbundl/pkgs/components/hero" // Register internal init()
	"bizbundl/pkgs/components/product_grid"
)

type PageBuilderModule struct {
	Service  *service.PageBuilderService
	Resolver *resolver.PageResolver
}

func Init(app *server.Server, catalogSvc *catalog_service.CatalogService, cartSvc *cart_service.CartService) *PageBuilderModule {
	svc := service.NewPageBuilderService(app.GetDB())

	// -- Atomic Component Registration --
	// Hero is registered via init() in pkg/components/hero
	// ProductGrid requires Service Injection
	product_grid.Register(catalogSvc)
	checkout.Register(cartSvc, catalogSvc)

	// Core Resolver
	res := resolver.NewPageResolver()

	// SeedDefaults removed from Global Init.
	// It should be called per-tenant when a shop is provisioned or accessed.

	return &PageBuilderModule{
		Service:  svc,
		Resolver: res,
	}
}
