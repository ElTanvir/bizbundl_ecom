package product_grid

import (
	"bizbundl/internal/storefront/catalog/service"
	"bizbundl/pkgs/components/registry"
	"bizbundl/pkgs/components/utils"

	// Explicit Imports
	"bizbundl/pkgs/components/product_grid/variants/carousel"
	"bizbundl/pkgs/components/product_grid/variants/grid"

	"github.com/a-h/templ"
)

// Register registers the ProductGrid component (Renderer + Resolver)
func Register(catalogSvc *service.CatalogService) {
	c := &registry.Component{
		Type:        "product_grid",
		Title:       "Product Grid",
		Description: "Displays a collection of products.",
		Category:    "Commerce",
		Variants:    make(map[string]registry.VariantDefinition),
		// Resolver will be attached later or initialized here
	}

	// Register Variants explicitly (Fixes Init Race Condition)
	c.Variants["grid"] = grid.Definition(catalogSvc)
	c.Variants["carousel"] = carousel.Definition(catalogSvc)

	// Dispatcher Resolver
	c.Resolver = NewResolver(catalogSvc)

	// Dispatcher Renderer
	c.Renderer = func(props map[string]interface{}) templ.Component {
		variant := utils.GetString(props, "Variant")
		// 1. Try Specific Variant Renderer
		if v, ok := c.Variants[variant]; ok && v.Renderer != nil {
			return v.Renderer(props)
		}
		// 2. Fallback to "Standard Grid" (grid)
		if v, ok := c.Variants["grid"]; ok && v.Renderer != nil {
			return v.Renderer(props)
		}

		return nil
	}

	registry.Register(c)
}
