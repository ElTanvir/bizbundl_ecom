package service_test

import (
	"context"
	"testing"

	"bizbundl/internal/storefront/auth/service"
	"bizbundl/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testutil.Cleanup(t)
	defer testutil.Cleanup(t)

	srv := testutil.SetupTestServer()
	svc := service.NewAuthService(srv.GetDB(), srv.GetTokenMaker())
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

	srv := testutil.SetupTestServer()
	svc := service.NewAuthService(srv.GetDB(), srv.GetTokenMaker())
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
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)

	// non-existent user
	_, _, err = svc.Login(ctx, "ghost@example.com", "any")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}
