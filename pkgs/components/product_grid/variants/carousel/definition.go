package carousel

import (
	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/storefront/catalog/service"
	"bizbundl/pkgs/components/registry"
	"bizbundl/pkgs/components/utils"
	"context"

	"github.com/a-h/templ"
)

func Definition(catalogSvc *service.CatalogService) registry.VariantDefinition {
	return registry.VariantDefinition{
		Name:        "carousel",
		Description: "Product Carousel",
		Props: map[string]registry.PropDefinition{
			"Title":       {Type: registry.TypeString, Default: "New Arrivals"},
			"Limit":       {Type: registry.TypeNumber, Default: 12},
			"CardVariant": {Type: registry.TypeString, Default: "standard"},
		},
		Renderer: func(props map[string]interface{}) templ.Component {
			return View(mapProps(props))
		},
		Resolver: &resolver{catalogSvc: catalogSvc},
	}
}

type resolver struct {
	catalogSvc *service.CatalogService
}

func (r *resolver) Resolve(ctx context.Context, section *registry.Section) error {
	limit := utils.GetInt(section.Props, "Limit")
	if limit == 0 {
		limit = 12 // Carousel Default
	}

	filter := utils.GetString(section.Props, "Filter")
	if filter == "" {
		filter = "new_arrivals" // Default for Carousel
	}

	var products []db.Product
	var err error

	switch filter {
	case "featured":
		products, err = r.catalogSvc.ListFeaturedProducts(ctx, int32(limit))
	case "new_arrivals":
		products, err = r.catalogSvc.ListNewArrivals(ctx, int32(limit))
	default:
		products, err = r.catalogSvc.ListProducts(ctx)
	}

	if err == nil {
		if len(products) > limit {
			products = products[:limit]
		}
		section.Props["Products"] = products
	}
	return err
}

func mapProps(props map[string]interface{}) Props {
	return Props{
		Title:       utils.GetString(props, "Title"),
		ViewAllLink: utils.GetString(props, "ViewAllLink"),
		Variant:     utils.GetString(props, "Variant"),
		CardVariant: utils.GetString(props, "CardVariant"),
		Products:    utils.GetProducList(props, "Products"),
	}
}
