package product_grid

import (
	"bizbundl/internal/storefront/catalog/service"
	"bizbundl/pkgs/components/registry"
	"context"
)

type Resolver struct {
	catalogService *service.CatalogService
}

func NewResolver(catalogService *service.CatalogService) *Resolver {
	return &Resolver{
		catalogService: catalogService,
	}
}

// Resolve delegates to the appropriate Variant Resolver.
func (r *Resolver) Resolve(ctx context.Context, section *registry.Section) error {
	variantName := "grid" // Default fallback
	if v, ok := section.Props["Variant"].(string); ok && v != "" {
		variantName = v
	}

	// 1. Get Component from Registry to access Variants
	// We assume "product_grid" is the type.
	comp, ok := registry.Get("product_grid")
	if !ok {
		return nil
	}

	// 2. Lookup Variant
	// First try the requested variant
	if variant, exists := comp.Variants[variantName]; exists {
		if variant.Resolver != nil {
			return variant.Resolver.Resolve(ctx, section)
		}
	}

	// 3. Fallback: If variant has no resolver, or variant not found (try default grid)
	if variantName != "grid" {
		if grid, exists := comp.Variants["grid"]; exists && grid.Resolver != nil {
			return grid.Resolver.Resolve(ctx, section)
		}
	}

	return nil
}
