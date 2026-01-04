package service

import (
	"context"
	"encoding/json"
	"fmt"

	db "bizbundl/internal/db/sqlc"

	"github.com/jackc/pgx/v5"
)

type Section struct {
	Type  string                 `json:"type"`
	Props map[string]interface{} `json:"props"`
}

type PageConfig struct {
	ID       string    `json:"id"`
	Route    string    `json:"route"`
	Title    string    `json:"title"`
	Sections []Section `json:"sections"`
}

type PageBuilderService struct {
	store db.DBStore
}

func NewPageBuilderService(store db.DBStore) *PageBuilderService {
	return &PageBuilderService{store: store}
}

func (s *PageBuilderService) GetPage(ctx context.Context, route string) (*PageConfig, error) {
	page, err := s.store.GetPageByRoute(ctx, route)
	if err != nil {
		return nil, err
	}

	var sections []Section
	if len(page.Sections) > 0 {
		if err := json.Unmarshal(page.Sections, &sections); err != nil {
			return nil, fmt.Errorf("failed to parse sections: %w", err)
		}
	}

	return &PageConfig{
		ID:       fmt.Sprintf("%x", page.ID.Bytes), // Simplified UUID string
		Route:    page.Route,
		Title:    page.Name,
		Sections: sections,
	}, nil
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
	defaultSections := []Section{
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
