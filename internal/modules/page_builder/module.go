package page_builder

import (
	cart_service "bizbundl/internal/modules/cart/service"
	catalog_service "bizbundl/internal/modules/catalog/service"
	"bizbundl/internal/modules/page_builder/resolver"
	"bizbundl/internal/modules/page_builder/service"
	"bizbundl/internal/server"
	"context"

	// Import Atomic Components
	"bizbundl/pkgs/components/checkout"
	_ "bizbundl/pkgs/components/hero" // Register internal init()
	"bizbundl/pkgs/components/product_grid"

	"github.com/rs/zerolog/log"
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

	// Seed Defaults
	if err := svc.SeedDefaults(context.Background()); err != nil {
		log.Error().Err(err).Msg("Failed to seed default pages")
	} else {
		log.Info().Msg("PageBuilder: Default pages seeded/verified")
	}

	return &PageBuilderModule{
		Service:  svc,
		Resolver: res,
	}
}
