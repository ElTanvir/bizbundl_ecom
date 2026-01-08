package service

import (
	"bizbundl/internal/store"
	"bizbundl/pkgs/components/registry"
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "bizbundl/internal/db/sqlc"

	"github.com/jackc/pgx/v5"
)

type PageConfig struct {
	ID       string             `json:"id"`
	Route    string             `json:"route"`
	Title    string             `json:"title"`
	Sections []registry.Section `json:"sections"`
}

type PageBuilderService struct {
	store db.DBStore
}

func NewPageBuilderService(store db.DBStore) *PageBuilderService {
	return &PageBuilderService{store: store}
}

func (s *PageBuilderService) GetPage(ctx context.Context, route string) (*PageConfig, error) {
	// 1. Cache Check
	cacheKey := "pb:page:" + route
	if val, ok := store.Get().Get(ctx, cacheKey); ok {
		return val.(*PageConfig), nil
	}

	// 2. DB Fetch
	page, err := s.store.GetPageByRoute(ctx, route)
	if err != nil {
		return nil, err
	}

	var sections []registry.Section
	if len(page.Sections) > 0 {
		if err := json.Unmarshal(page.Sections, &sections); err != nil {
			return nil, fmt.Errorf("failed to parse sections: %w", err)
		}
	}

	cfg := &PageConfig{
		ID:       fmt.Sprintf("%x", page.ID.Bytes),
		Route:    page.Route,
		Title:    page.Name,
		Sections: sections,
	}

	// 3. Cache Set (Pointer) - 24 Hours
	store.Get().Set(ctx, cacheKey, cfg, 24*time.Hour)

	return cfg, nil
}

// ValidatePage checks if the page structure adheres to registry constraints (e.g. AllowedChildren)
func (s *PageBuilderService) ValidatePage(sections []registry.Section) error {
	for i, section := range sections {
		// 1. Check if Component Exists
		comp, exists := registry.Get(section.Type)
		if !exists {
			return fmt.Errorf("section %d: component type '%s' not found", i, section.Type)
		}

		// 2. Check Children Constraints
		if childrenRaw, ok := section.Props["children"]; ok {
			// Try to marshal/unmarshal to []Section to be safe with interface{}
			// Or just assume it's []interface{} and inspect "type"
			// For robustness, let's assume standard JSON unmarshal would give []interface{} map[string]interface{}

			// Simple check: If component doesn't allow children but has them -> Error (optional strictness)
			// But main check is: If it HAS allowed_children, are these children VALID?

			if len(comp.AllowedChildren) > 0 {
				var children []registry.Section
				// Helper to convert props["children"] to []Section
				bytes, _ := json.Marshal(childrenRaw)
				json.Unmarshal(bytes, &children)

				for _, child := range children {
					isAllowed := false
					for _, allowed := range comp.AllowedChildren {
						if allowed == child.Type {
							isAllowed = true
							break
						}
					}
					if !isAllowed {
						return fmt.Errorf("component '%s' does not allow child '%s'. Allowed: %v", section.Type, child.Type, comp.AllowedChildren)
					}

					// Recurse
					if err := s.ValidatePage([]registry.Section{child}); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (s *PageBuilderService) SeedDefaults(ctx context.Context) error {
	// Check if Home exists
	_, err := s.store.GetPageByRoute(ctx, "/")
	if err == nil {
		return nil // Already exists
	}

	if err != pgx.ErrNoRows {
		return err // Real error
	}

	// Create Default Home
	defaultSections := []registry.Section{
		{
			Type: "hero",
			Props: map[string]interface{}{
				"Title":      "Welcome to BizBundl",
				"Subtitle":   "Premium Digital Assets for your Business",
				"ButtonText": "Browse Catalog",
				"ButtonLink": "/product/demo-pro", // Example
				"Align":      "center",
			},
		},
		{
			Type: "product_grid",
			Props: map[string]interface{}{
				"Title": "Featured Products",
				"Limit": 4,
			},
		},
	}

	sectionsJSON, _ := json.Marshal(defaultSections)
	isPublished := true

	_, err = s.store.CreatePage(ctx, db.CreatePageParams{
		Route:       "/",
		Name:        "Home",
		Sections:    sectionsJSON,
		IsPublished: &isPublished,
	})

	return err
}
