package resolver

import (
	pb "bizbundl/internal/modules/page_builder/service"
	"context"
)

// ComponentResolver defines the contract for resolving data for a specific component type.
type ComponentResolver interface {
	Resolve(ctx context.Context, section *pb.Section) error
}

// Registry holds the mapping of component types to their resolvers.
type Registry struct {
	resolvers map[string]ComponentResolver
}

func NewRegistry() *Registry {
	return &Registry{
		resolvers: make(map[string]ComponentResolver),
	}
}

func (r *Registry) Register(componentType string, resolver ComponentResolver) {
	r.resolvers[componentType] = resolver
}

func (r *Registry) Get(componentType string) (ComponentResolver, bool) {
	res, ok := r.resolvers[componentType]
	return res, ok
}

// Resolve delegates resolving to the registered component resolver.
func (r *Registry) Resolve(ctx context.Context, section *pb.Section) error {
	resolver, ok := r.Get(section.Type)
	if !ok {
		// No resolver needed for this type (e.g., static content like rich_text)
		return nil
	}
	return resolver.Resolve(ctx, section)
}
