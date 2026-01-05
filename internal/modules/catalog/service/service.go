package service

import (
	"context"
	"fmt"

	"strings"

	db "bizbundl/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type CatalogService struct {
	store db.DBStore
}

func NewCatalogService(store db.DBStore) *CatalogService {
	return &CatalogService{store: store}
}

// -- Categories --

func (s *CatalogService) CreateCategory(ctx context.Context, name string, parentID pgtype.UUID) (db.Category, error) {
	slug := makeSlug(name)
	return s.store.CreateCategory(ctx, db.CreateCategoryParams{
		Name:     name,
		Slug:     slug,
		ParentID: parentID,
		IsActive: boolPtr(true),
	})
}

func (s *CatalogService) GetCategory(ctx context.Context, id pgtype.UUID) (db.Category, error) {
	return s.store.GetCategory(ctx, id)
}

func (s *CatalogService) ListCategories(ctx context.Context) ([]db.Category, error) {
	return s.store.ListCategories(ctx)
}

// -- Products --

type CreateProductParams struct {
	Title       string
	Description string
	BasePrice   float64
	IsDigital   bool
	FilePath    string
	CategoryID  pgtype.UUID
	IsFeatured  bool
}

func (s *CatalogService) CreateProduct(ctx context.Context, p CreateProductParams) (db.Product, error) {
	slug := makeSlug(p.Title)

	priceNumeric := pgtype.Numeric{}
	err := priceNumeric.Scan(fmt.Sprintf("%f", p.BasePrice))
	if err != nil {
		return db.Product{}, fmt.Errorf("invalid price: %v", err)
	}

	return s.store.CreateProduct(ctx, db.CreateProductParams{
		Title:       p.Title,
		Slug:        slug,
		Description: strPtr(p.Description),
		BasePrice:   priceNumeric,
		IsDigital:   boolPtr(p.IsDigital),
		FilePath:    strPtr(p.FilePath),
		CategoryID:  p.CategoryID,
		IsActive:    boolPtr(true),
		IsFeatured:  boolPtr(p.IsFeatured),
	})
}

func (s *CatalogService) GetProduct(ctx context.Context, id pgtype.UUID) (db.Product, error) {
	return s.store.GetProduct(ctx, id)
}

func (s *CatalogService) GetProductBySlug(ctx context.Context, slug string) (db.Product, error) {
	return s.store.GetProductBySlug(ctx, slug)
}

func (s *CatalogService) ListProducts(ctx context.Context) ([]db.Product, error) {
	return s.store.ListProducts(ctx)
}

func (s *CatalogService) ListFeaturedProducts(ctx context.Context, limit int32) ([]db.Product, error) {
	return s.store.ListFeaturedProducts(ctx, limit)
}

func (s *CatalogService) ListNewArrivals(ctx context.Context, limit int32) ([]db.Product, error) {
	return s.store.ListNewArrivals(ctx, limit)
}

// -- Utilities --

func makeSlug(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}

func boolPtr(b bool) *bool {
	return &b
}
func strPtr(s string) *string {
	return &s
}
