package resolver

import (
	pb "bizbundl/internal/modules/page_builder/service"
	"bizbundl/pkgs/components/registry"
	"context"
	"fmt"
)

type PageResolver struct{}

func NewPageResolver() *PageResolver {
	return &PageResolver{}
}

// Resolve iterates through the page sections and delegates to the Registry.
func (r *PageResolver) Resolve(ctx context.Context, page *pb.PageConfig) error {
	for i := range page.Sections {
		// Look up component in Global Registry
		if comp, ok := registry.Get(page.Sections[i].Type); ok && comp.Resolver != nil {
			if err := comp.Resolver.Resolve(ctx, &page.Sections[i]); err != nil {
				fmt.Printf("Error resolving section %s: %v\n", page.Sections[i].Type, err)
			}
		}
	}
	return nil
}
