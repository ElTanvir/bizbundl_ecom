package main

import (
	"bizbundl/internal/config"
	"bizbundl/util"
	"fmt"
	"log"
)

// Tenants Registry (MVP: Hardcoded / Config-based)
// In production, fetch this from centralized Redis or a 'public.shops' table
var activeTenants = []string{
	"shop_1", // Default/First shop
	"shop_2", // Test shop
}

func main() {
	fmt.Println("🚀 Starting Multi-Schema Migration Worker...")

	// 1. Load Config (DB Auth)
	cfg := config.Load()

	// Base URL without search_path
	// DBSourceURL typically looks like: postgres://user:pass@host:5432/dbname?sslmode=disable
	baseURL := cfg.DBSourceURL()

	// 2. Migration Directories
	platformMigDir := "internal/db/migration/platform"
	tenantMigDir := "internal/db/migration/tenant"

	if cfg.InDocker == "true" {
		platformMigDir = "/app/internal/db/migration/platform"
		tenantMigDir = "/app/internal/db/migration/tenant"
	}

	// 3. Migrate 'public' Schema (Shared Infrastructure)
	fmt.Println(">> Migrating 'public' schema...")
	// For public, we might want to explicity set search_path=public, or default
	// Let's force it to be safe
	publicURL := fmt.Sprintf("%s&search_path=public", baseURL)
	if err := util.RunMigrations(publicURL, platformMigDir); err != nil {
		log.Fatalf("❌ Public Migration Failed: %v", err)
	}
	fmt.Println("✅ Public Schema Up-to-Date.")

	// 4. Iterate Tenants
	for _, tenantID := range activeTenants {
		fmt.Printf(">> Migrating Tenant: %s...\n", tenantID)

		// Construct DSN with search_path
		// Note: golang-migrate postgres driver respects search_path query param
		tenantURL := fmt.Sprintf("%s&search_path=%s,public", baseURL, tenantID)

		// Create Schema if not exists (This worker creates schemas for now!)
		// Ideally, the App creates the schema on signup.
		// But for 'migrations', we assume schema exists.
		// However, standard PG driver doesn't CREATE SCHEMA automatically.
		// For MVP simplicity, we assume schema exists OR fail.
		// Detailed logic would connect to public, run 'CREATE SCHEMA IF NOT EXISTS', then migrate.
		// Let's just RunMigrations. If schema missing, it fails.

		if err := util.RunMigrations(tenantURL, tenantMigDir); err != nil {
			log.Printf("⚠️  Migration Failed for %s: %v", tenantID, err)
			// Continue to next tenant; don't potentially crash entire fleet for one bad tenant state
			continue
		}
		fmt.Printf("✅ Tenant %s Migrated.\n", tenantID)
	}

	fmt.Println("🏁 All Migrations Completed.")
}
