package testutil

import (
	"context"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"bizbundl/internal/config"
	db "bizbundl/internal/db/sqlc"
	"bizbundl/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore db.DBStore
var testPool *pgxpool.Pool

func SetupTestDB() db.DBStore {
	if testStore != nil {
		return testStore
	}

	// Find Project Root
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	root := filepath.Join(basepath, "../..")

	cfg := config.Load()

	// Run Migrations (Reset for clean test state)
	migrationPath := filepath.Join(root, "internal/db/migration")
	// cfg.DBSourceURL() is available
	if err := util.RunMigrationReset(cfg.DBSourceURL(), migrationPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	connPool, err := pgxpool.New(context.Background(), cfg.DBSource())
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testPool = connPool

	testStore = db.NewStore(connPool)
	return testStore
}

// CleanupTables truncates tables to ensure clean state between tests
func Cleanup(t *testing.T) {
	if testPool == nil {
		return
	}

	// Order matters due to FKs
	tables := []string{
		"cart_items", "carts",
		"order_items", "orders",
		"sessions",
		"product_variants", "products", "categories",
		"users",
	}

	q := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE;"
	_, err := testPool.Exec(context.Background(), q)
	if err != nil {
		t.Logf("Failed to truncate tables: %v", err)
	}
}

func RandomEmail() string {
	return RandomString(10) + "@example.com"
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
