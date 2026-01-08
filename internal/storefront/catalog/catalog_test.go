package catalog_test

import (
	"context"
	"testing"

	"bizbundl/internal/storefront/catalog/service"
	"bizbundl/internal/testutil"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	srv := testutil.SetupTestServer()
	store := srv.GetDB()
	svc := service.NewCatalogService(store)
	ctx := context.Background()

	cat, err := svc.CreateCategory(ctx, "Electronics", pgtype.UUID{})
	require.NoError(t, err)
	assert.Equal(t, "Electronics", cat.Name)
	assert.Equal(t, "electronics", cat.Slug)
	assert.True(t, *cat.IsActive)
}

func TestCreateProduct(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	srv := testutil.SetupTestServer()
	store := srv.GetDB()
	svc := service.NewCatalogService(store)
	ctx := context.Background()

	// Dependency: Category
	cat, err := svc.CreateCategory(ctx, "Laptops", pgtype.UUID{})
	assert.NoError(t, err)

	p, err := svc.CreateProduct(ctx, service.CreateProductParams{
		Title:      "MacBook Pro M4",
		BasePrice:  1999.00,
		IsDigital:  false,
		FilePath:   "",
		CategoryID: cat.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, "macbook-pro-m4", p.Slug)
	assert.Equal(t, cat.ID, p.CategoryID)
}

func TestGetProduct(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	srv := testutil.SetupTestServer()
	store := srv.GetDB()
	svc := service.NewCatalogService(store)
	ctx := context.Background()

	cat, _ := svc.CreateCategory(ctx, "Phones", pgtype.UUID{})
	pCreated, _ := svc.CreateProduct(ctx, service.CreateProductParams{
		Title:      "iPhone 16",
		BasePrice:  999.00,
		CategoryID: cat.ID,
	})

	// Retrieve
	p, err := svc.GetProduct(ctx, pCreated.ID)
	assert.NoError(t, err)
	assert.Equal(t, pCreated.Title, p.Title)

	// Fail
	_, err = svc.GetProduct(ctx, pgtype.UUID{Valid: false}) // Invalid UUID
	assert.Error(t, err)
}

func TestListCategories(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	srv := testutil.SetupTestServer()
	store := srv.GetDB()
	svc := service.NewCatalogService(store)
	ctx := context.Background()

	_, _ = svc.CreateCategory(ctx, "A", pgtype.UUID{})
	_, _ = svc.CreateCategory(ctx, "B", pgtype.UUID{})

	cats, err := svc.ListCategories(ctx)
	assert.NoError(t, err)
	assert.Len(t, cats, 2)
}

func TestListProducts(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	srv := testutil.SetupTestServer()
	store := srv.GetDB()
	svc := service.NewCatalogService(store)
	ctx := context.Background()

	cat, _ := svc.CreateCategory(ctx, "C", pgtype.UUID{})
	_, _ = svc.CreateProduct(ctx, service.CreateProductParams{Title: "P1", CategoryID: cat.ID})

	items, err := svc.ListProducts(ctx)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
}
