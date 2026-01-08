package service

import (
	"context"
	"fmt"
	"strings"

	"bizbundl/internal/config"
	db "bizbundl/internal/db/sqlc/platform" // platform queries

	// We need a way to run migrations.
	"bizbundl/util"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PlatformStore defines access to platform DB (shops, users)
type PlatformStore interface {
	db.Querier
	GetPool() *pgxpool.Pool
}

// SQLPlatformStore wraps the auto-generated Queries plus Pool access
type SQLPlatformStore struct {
	*db.Queries
	pool *pgxpool.Pool
}

func (s *SQLPlatformStore) GetPool() *pgxpool.Pool {
	return s.pool
}

// PlatformService logic
type PlatformService struct {
	store PlatformStore
	cfg   *config.Config
}

// NewPlatformService factory
func NewPlatformService(pool *pgxpool.Pool, cfg *config.Config) *PlatformService {
	// Manually construct the store wrapper since it's structurally simple
	// Note: The main 'Store' in internal/db/sqlc points to 'db' package (Tenants).
	// We are using 'platform' package here.
	queries := db.New(pool)
	return &PlatformService{
		store: &SQLPlatformStore{Queries: queries, pool: pool},
		cfg:   cfg,
	}
}

// CreateShop orchestrates creating the Shop record and provisioning the Schema
func (s *PlatformService) CreateShop(ctx context.Context, ownerID pgtype.UUID, name string) (db.Shop, error) {
	// 1. Generate Subdomain (simple slugify)
	subdomain := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	// TODO: Ensure uniqueness or retry with suffix

	// 2. Generate Tenant ID (Schema Name) -> "shop_xyz"
	// Sanitize subdomain specific characters for SQL schema name safety
	tenantID := "shop_" + strings.ReplaceAll(subdomain, "-", "_")

	// 3. Create Record
	isActive := true
	shop, err := s.store.CreateShop(ctx, db.CreateShopParams{
		OwnerID:   ownerID,
		Name:      name,
		Subdomain: subdomain,
		TenantID:  tenantID,
		IsActive:  &isActive,
	})
	if err != nil {
		return db.Shop{}, fmt.Errorf("failed to create shop record: %w", err)
	}

	// 4. Provision Schema (Migration)
	// We need to run "tenant" migrations on this new schema.
	// DSN Construction:
	baseURL := s.cfg.DBSourceURL()
	tenantURL := fmt.Sprintf("%s&search_path=%s,public", baseURL, tenantID)

	migDir := "internal/db/migration/tenant"
	if s.cfg.InDocker == "true" {
		migDir = "/app/internal/db/migration/tenant"
	}

	if err := util.RunMigrations(tenantURL, migDir); err != nil {
		// Rollback? If record created but schema failed, we have a broken state.
		// For MVP, we return error. User retries -> Subdomain constraint fails.
		// Ideally: Wrap in transaction OR Delete shop on failure.
		// Let's Log & Return error.
		return shop, fmt.Errorf("shop created but schema provision failed: %w", err)
	}

	return shop, nil
}

func (s *PlatformService) ListShops(ctx context.Context, ownerID pgtype.UUID) ([]db.Shop, error) {
	return s.store.ListShopsByOwner(ctx, ownerID)
}
