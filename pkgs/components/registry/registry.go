package registry

import (
	"context"

	"github.com/a-h/templ"
)

// -- Definition --

// Section represents an instance of a component type with specific props
type Section struct {
	Type  string                 `json:"type"`
	Props map[string]interface{} `json:"props"`
}

type Resolvable interface {
	Resolve(ctx context.Context, section *Section) error
}

type RenderFunc func(props map[string]interface{}) templ.Component

// PropType defines the expected type of a prop
type PropType string

const (
	TypeString  PropType = "string"
	TypeNumber  PropType = "number"
	TypeBoolean PropType = "boolean"
	TypeImage   PropType = "image"
	TypeColor   PropType = "color"
)

// PropDefinition describes a single property for the editor/AI
type PropDefinition struct {
	Type        PropType    `json:"type"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// VariantDefinition defines the schema for a specific visual variant
type VariantDefinition struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description,omitempty"`
	Props       map[string]PropDefinition `json:"props"` // Schema specific to this variant
	Renderer    RenderFunc                `json:"-"`     // Specific Renderer
	Resolver    Resolvable                `json:"-"`     // Specific Resolver
}

type Component struct {
	Type     string
	Renderer RenderFunc
	Resolver Resolvable // Optional

	// Advanced Builder Config
	Title           string                       `json:"title"`
	Description     string                       `json:"description"`
	Category        string                       `json:"category"`                   // e.g. "Layout", "Content", "Commerce"
	AllowedChildren []string                     `json:"allowed_children,omitempty"` // Whitelist of component types (Slots)
	Variants        map[string]VariantDefinition `json:"variants"`                   // Variant-specific schemas
}

// -- Registry (Singleton) --

var components = make(map[string]*Component)

func Register(c *Component) {
	if _, exists := components[c.Type]; exists {
		panic("registry: duplicate component registration for type: " + c.Type)
	}

	// Initialize Variants map if nil
	if c.Variants == nil {
		c.Variants = make(map[string]VariantDefinition)
	}

	components[c.Type] = c
}

// RegisterVariant allows adding a variant to an existing component.
// This enables splitting variant definitions into separate files/packages.
func RegisterVariant(componentType string, v VariantDefinition) {
	comp, exists := components[componentType]
	if !exists {
		panic("registry: cannot register variant for unknown component type: " + componentType)
	}

	if _, exists := comp.Variants[v.Name]; exists {
		panic("registry: duplicate variant '" + v.Name + "' for component: " + componentType)
	}

	comp.Variants[v.Name] = v
}

func Get(compType string) (*Component, bool) {
	c, ok := components[compType]
	return c, ok
}

// -- Helpers --

func GetAllHelpers() map[string]*Component {
	return components
}

// GetJSONSchema returns the full registry schema for MCP/AI agents
func GetJSONSchema() map[string]*Component {
	return components
}
