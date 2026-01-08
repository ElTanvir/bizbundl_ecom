package cart_test

import (
	"context"
	"testing"
	"time"

	db "bizbundl/internal/db/sqlc"
	authservice "bizbundl/internal/storefront/auth/service"
	"bizbundl/internal/storefront/cart/service"
	catalogservice "bizbundl/internal/storefront/catalog/service"
	"bizbundl/internal/testutil"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddToCart(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	srv := testutil.SetupTestServer()
	store := srv.GetDB()
	cartSvc := service.NewCartService(store)
	catalogSvc := catalogservice.NewCatalogService(store)
	ctx := context.Background()

	// Setup Product
	cat, _ := catalogSvc.CreateCategory(ctx, "Test Cat", pgtype.UUID{})
	p, err := catalogSvc.CreateProduct(ctx, catalogservice.CreateProductParams{
		Title: "Item 1", BasePrice: 100.0, CategoryID: cat.ID,
	})
	require.NoError(t, err)

	// Helper to create guest session
	createGuestSession := func() pgtype.UUID {
		id := uuid.New()
		token := id.String() // simple token
		expires := time.Now().Add(1 * time.Hour)

		s, err := store.CreateSession(ctx, db.CreateSessionParams{
			Token:     token,
			UserID:    pgtype.UUID{}, // Null
			ExpiresAt: pgtype.Timestamptz{Time: expires, Valid: true},
		})
		require.NoError(t, err)
		return s.ID
	}

	sessionID := createGuestSession()

	// Add Item (1)
	item, err := cartSvc.AddToCart(ctx, sessionID, pgtype.UUID{}, p.ID, pgtype.UUID{}, 1)
	require.NoError(t, err)
	assert.Equal(t, int32(1), item.Quantity)

	// Add Same Item (Upsert +2)
	item, err = cartSvc.AddToCart(ctx, sessionID, pgtype.UUID{}, p.ID, pgtype.UUID{}, 2)
	require.NoError(t, err)
	assert.Equal(t, int32(3), item.Quantity)

	// Check Cart Created
	c, err := cartSvc.GetOrCreateCart(ctx, sessionID, pgtype.UUID{})
	require.NoError(t, err)
	assert.Equal(t, sessionID, c.SessionID)
}

func TestMergeCarts(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	srv := testutil.SetupTestServer()
	store := srv.GetDB()
	cartSvc := service.NewCartService(store)
	catalogSvc := catalogservice.NewCatalogService(store)
	authSvc := authservice.NewAuthService(store, srv.GetTokenMaker())
	ctx := context.Background()

	// Product
	cat, _ := catalogSvc.CreateCategory(ctx, "Test Cat 2", pgtype.UUID{})
	p, _ := catalogSvc.CreateProduct(ctx, catalogservice.CreateProductParams{Title: "P1", BasePrice: 10.0, CategoryID: cat.ID})

	// Guest Session (Valid)
	guestSessionID := uuid.New()
	guestToken := guestSessionID.String()
	sess, err := store.CreateSession(ctx, db.CreateSessionParams{
		Token:     guestToken,
		UserID:    pgtype.UUID{},
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true},
	})
	require.NoError(t, err)

	// Create Cart with Session ID (which is sess.ID)
	_, err = cartSvc.AddToCart(ctx, sess.ID, pgtype.UUID{}, p.ID, pgtype.UUID{}, 5)
	require.NoError(t, err)

	// User Login
	email := testutil.RandomEmail()
	_, _ = authSvc.Register(ctx, email, "pass", "User", "123")
	_, user, _ := authSvc.Login(ctx, email, "pass")

	// Merge
	err = cartSvc.MergeCarts(ctx, sess.ID, user.ID)
	require.NoError(t, err)

	// Verify User has items
	userCart, err := cartSvc.GetOrCreateCart(ctx, pgtype.UUID{}, user.ID) // By ID
	require.NoError(t, err)

	items, err := cartSvc.GetCartItems(ctx, userCart.ID)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, int32(5), items[0].Quantity)
}
