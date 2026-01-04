package auth_test

import (
	"context"
	"testing"
	"time"

	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/modules/auth"
	"bizbundl/internal/testutil"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	store := testutil.SetupTestDB()
	svc := auth.NewAuthService(store)
	ctx := context.Background()

	email := testutil.RandomEmail()

	// Good Case
	user, err := svc.Register(ctx, email, "password123", "Test User", "1234567890")
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.NotEmpty(t, user.PasswordHash)

	// Duplicate Email
	_, err = svc.Register(ctx, email, "newpass", "New User", "0000000000")
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	store := testutil.SetupTestDB()
	svc := auth.NewAuthService(store)
	ctx := context.Background()

	email := testutil.RandomEmail()
	_, _ = svc.Register(ctx, email, "securepass", "Login User", "111")

	// correct login
	token, user, err := svc.Login(ctx, email, "securepass")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, email, user.Email)

	// wrong password
	_, _, err = svc.Login(ctx, email, "wrongpass")
	assert.ErrorIs(t, err, auth.ErrInvalidCredentials)

	// non-existent user
	_, _, err = svc.Login(ctx, "ghost@example.com", "any")
	assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
}

func TestVerifySession(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	store := testutil.SetupTestDB()
	svc := auth.NewAuthService(store)
	ctx := context.Background()

	email := testutil.RandomEmail()
	_, _ = svc.Register(ctx, email, "pass", "Session User", "222")
	token, user, err := svc.Login(ctx, email, "pass")
	require.NoError(t, err)

	// Valid
	u, err := svc.VerifySession(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, email, u.Email)
	assert.Equal(t, user.ID, u.ID)

	// Invalid Token
	_, err = svc.VerifySession(ctx, "fake-token")
	assert.Error(t, err)

	// Expired Token (Manually insert expired session into Real DB)
	expiredToken := "expired-uuid"
	expiresAt := time.Now().Add(-1 * time.Hour)

	// Direct DB insert helper or just ignoring encapsulation for test setup?
	// We can use store directly.
	_, err = store.CreateSession(ctx, db.CreateSessionParams{
		Token:     expiredToken,
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	assert.NoError(t, err)

	_, err = svc.VerifySession(ctx, expiredToken)
	assert.Error(t, err) // Should error and delete

	// Verify deletion
	_, err = store.GetSession(ctx, expiredToken)
	assert.Error(t, err)
}

func TestLogout(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	store := testutil.SetupTestDB()
	svc := auth.NewAuthService(store)
	ctx := context.Background()

	email := testutil.RandomEmail()
	_, _ = svc.Register(ctx, email, "pass", "Logout User", "333")
	tok, _, _ := svc.Login(ctx, email, "pass")

	err := svc.Logout(ctx, tok)
	assert.NoError(t, err)

	_, err = svc.VerifySession(ctx, tok)
	assert.Error(t, err)
}
