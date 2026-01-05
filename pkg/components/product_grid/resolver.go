package product_grid

import (
	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/modules/catalog/service"
	pb "bizbundl/internal/modules/page_builder/service"
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

func (r *Resolver) Resolve(ctx context.Context, section *pb.Section) error {
	// 1. Extract Props
	limit := 4
	if l, ok := section.Props["Limit"].(float64); ok {
		limit = int(l)
	}

	filter, _ := section.Props["Filter"].(string)

	// 2. Fetch Data
	var products []db.Product
	var err error

	switch filter {
	case "featured":
		products, err = r.catalogService.ListFeaturedProducts(ctx, int32(limit))
	case "new_arrivals":
		products, err = r.catalogService.ListNewArrivals(ctx, int32(limit))
	default:
		// Fallback
		products, err = r.catalogService.ListProducts(ctx)
	}

	if err != nil {
		return err
	}

	// 3. Apply Limit
	if len(products) > limit {
		section.Props["Products"] = products[:limit]
	} else {
		section.Props["Products"] = products
	}

	return nil
}
