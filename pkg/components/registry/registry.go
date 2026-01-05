package registry

import (
	pb "bizbundl/internal/modules/page_builder/service"
	"context"

	"github.com/a-h/templ"
)

// -- Definition --

type Resolvable interface {
	Resolve(ctx context.Context, section *pb.Section) error
}

type RenderFunc func(props map[string]interface{}) templ.Component

type Component struct {
	Type     string
	Renderer RenderFunc
	Resolver Resolvable // Optional
}

// -- Registry (Singleton) --

var components = make(map[string]*Component)

func Register(c *Component) {
	components[c.Type] = c
}

func Get(compType string) (*Component, bool) {
	c, ok := components[compType]
	return c, ok
}

// -- Helpers --

func GetAllHelpers() map[string]*Component {
	return components
}
