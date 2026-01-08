package main

import (
	"context"
	"fmt"
	"log"

	"bizbundl/internal/config"
	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/modules/catalog/service"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Connect DB
	cfg := config.Load()
	dbURL := cfg.DBSourceURL()

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Unable to parse DATABASE_URL:", err)
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	catalogSvc := service.NewCatalogService(store)
	ctx := context.Background()

	fmt.Println("🌱 Starting Seeding...")

	// 2. Seed Categories
	cats := []string{"SaaS Kits", "UI Templates", "E-Books", "Graphics"}
	var catIDs []pgtype.UUID

	for _, c := range cats {
		// Simple check via slug would be better, but for MVP we just create and ignore err if dup slug
		// Service CreateCategory generates slug.
		cat, err := catalogSvc.CreateCategory(ctx, c, pgtype.UUID{})
		if err != nil {
			// Assume duplicate, try to find it? Or just continue
			fmt.Printf("Category '%s' skipped or failed: %v\n", c, err)
			continue
		}
		fmt.Printf("Created Category: %s\n", c)
		catIDs = append(catIDs, cat.ID)
	}

	// Fetch all cats to ensure we have IDs even if create failed
	allCats, _ := catalogSvc.ListCategories(ctx)
	if len(allCats) == 0 {
		log.Fatal("No categories available to seed products")
	}

	// 3. Seed Products
	products := []service.CreateProductParams{
		{
			Title:       "E-Commerce Starter Kit",
			Description: "A full stack e-commerce solution in Go and Templ.",
			BasePrice:   199.99,
			IsDigital:   true,
			IsFeatured:  true,
		},
		{
			Title:       "Marketing Dashboard UI",
			Description: "High conversion dashboard template.",
			BasePrice:   49.00,
			IsDigital:   true,
			IsFeatured:  true,
		},
		{
			Title:       "Golang Mastery E-Book",
			Description: "Master Go in 30 days.",
			BasePrice:   29.99,
			IsDigital:   true,
			IsFeatured:  false,
		},
		{
			Title:       "Icon Pack: Neon",
			Description: "500+ Neon icons for dark mode.",
			BasePrice:   15.00,
			IsDigital:   true,
			IsFeatured:  false,
		},
		{
			Title:       "SaaS Boilerplate Pro",
			Description: "Enterprise grade boilerplate.",
			BasePrice:   299.00,
			IsDigital:   true,
			IsFeatured:  true,
		},
		{
			Title:       "Figma Wireframe Kit",
			Description: "Rapid prototyping kit.",
			BasePrice:   0.00, // Freebie?
			IsDigital:   true,
			IsFeatured:  true,
		},
	}

	for i, p := range products {
		// Assign random category
		p.CategoryID = allCats[i%len(allCats)].ID
		p.FilePath = "/downloads/dummy.zip"

		_, err := catalogSvc.CreateProduct(ctx, p)
		if err != nil {
			fmt.Printf("Product '%s' skipped: %v\n", p.Title, err)
		} else {
			fmt.Printf("Created Product: %s\n", p.Title)
		}
	}

	fmt.Println("✅ Seeding Complete!")
}
