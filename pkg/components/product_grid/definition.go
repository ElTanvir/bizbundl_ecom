package product_grid

import (
	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/modules/catalog/service"
	"bizbundl/pkg/components/registry"

	"github.com/a-h/templ"
)

// Register registers the ProductGrid component (Renderer + Resolver)
func Register(catalogSvc *service.CatalogService) {
	registry.Register(&registry.Component{
		Type:     "product_grid",
		Resolver: NewResolver(catalogSvc),
		Renderer: func(props map[string]interface{}) templ.Component {
			return View(mapProps(props))
		},
	})
}

func mapProps(props map[string]interface{}) Props {
	p := Props{
		Title:       getString(props, "Title"),
		ViewAllLink: getString(props, "ViewAllLink"),
		Variant:     getString(props, "Variant"),
	}

	// Extract Products
	if products, ok := props["Products"].([]db.Product); ok {
		p.Products = products
	}
	return p
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
