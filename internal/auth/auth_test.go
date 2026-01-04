package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"bizbundl/internal/auth"
	db "bizbundl/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

// --- Manual Mock ---

type MockStore struct {
	db.DBStore
	users    map[string]db.User // email -> user
	sessions map[string]db.Session
}

func NewMockStore() *MockStore {
	return &MockStore{
		users:    make(map[string]db.User),
		sessions: make(map[string]db.Session),
	}
}

func (m *MockStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	if _, exists := m.users[arg.Email]; exists {
		return db.User{}, errors.New("unique_violation")
	}

	// Convert [16]byte UUID to pgtype.UUID strictly
	pgID := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	user := db.User{
		ID:           pgID,
		Email:        arg.Email,
		PasswordHash: arg.PasswordHash,
		FullName:     arg.FullName,
		Role:         arg.Role,
	}
	if arg.Phone != nil {
		user.Phone = arg.Phone
	}
	m.users[arg.Email] = user
	return user, nil
}

func (m *MockStore) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	u, ok := m.users[email]
	if !ok {
		return db.User{}, errors.New("no rows")
	}
	return u, nil
}

func (m *MockStore) GetUserById(ctx context.Context, id pgtype.UUID) (db.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return db.User{}, errors.New("no rows")
}

func (m *MockStore) CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
	s := db.Session{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Token:     arg.Token,
		UserID:    arg.UserID,
		ExpiresAt: arg.ExpiresAt,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	m.sessions[arg.Token] = s
	return s, nil
}

func (m *MockStore) GetSession(ctx context.Context, token string) (db.Session, error) {
	s, ok := m.sessions[token]
	if !ok {
		return db.Session{}, errors.New("no rows")
	}
	return s, nil
}

func (m *MockStore) DeleteSession(ctx context.Context, token string) error {
	delete(m.sessions, token)
	return nil
}

// --- Tests ---

func TestRegister(t *testing.T) {
	store := NewMockStore()
	svc := auth.NewAuthService(store)
	ctx := context.Background()

	// Good Case
	user, err := svc.Register(ctx, "test@example.com", "password123", "Test User", "1234567890")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, user.PasswordHash)

	// Duplicate Email
	_, err = svc.Register(ctx, "test@example.com", "newpass", "New User", "0000000000")
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	store := NewMockStore()
	svc := auth.NewAuthService(store)
	ctx := context.Background()

	// Setup User
	_, _ = svc.Register(ctx, "login@example.com", "securepass", "Login User", "111")

	// correct login
	token, user, err := svc.Login(ctx, "login@example.com", "securepass")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "login@example.com", user.Email)

	// wrong password
	_, _, err = svc.Login(ctx, "login@example.com", "wrongpass")
	assert.ErrorIs(t, err, auth.ErrInvalidCredentials)

	// non-existent user
	_, _, err = svc.Login(ctx, "ghost@example.com", "any")
	assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
}

func TestVerifySession(t *testing.T) {
	store := NewMockStore()
	svc := auth.NewAuthService(store)
	ctx := context.Background()

	// Setup
	_, _ = svc.Register(ctx, "session@example.com", "pass", "Session User", "222")
	token, _, _ := svc.Login(ctx, "session@example.com", "pass")

	// Valid
	user, err := svc.VerifySession(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, "session@example.com", user.Email)

	// Invalid Token
	_, err = svc.VerifySession(ctx, "fake-token")
	assert.Error(t, err)

	// Expired Token (Manually manipulate mock)
	expiredToken := "expired-uuid"
	store.sessions[expiredToken] = db.Session{
		Token:     expiredToken,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(-1 * time.Hour), Valid: true},
		UserID:    user.ID,
	}

	_, err = svc.VerifySession(ctx, expiredToken)
	assert.Error(t, err) // Should error and delete
	_, exists := store.sessions[expiredToken]
	assert.False(t, exists, "Expired session should be deleted")
}

func TestLogout(t *testing.T) {
	store := NewMockStore()
	svc := auth.NewAuthService(store)
	ctx := context.Background()

	_, _, _ = svc.Login(ctx, "session@example.com", "pass") // assuming user from prev test doesn't exist in new store instance
	// Actually each test has new store. Need to register first.
	_, _ = svc.Register(ctx, "logout@example.com", "pass", "Logout User", "333")
	tok, _, _ := svc.Login(ctx, "logout@example.com", "pass")

	err := svc.Logout(ctx, tok)
	assert.NoError(t, err)

	_, err = svc.VerifySession(ctx, tok)
	assert.Error(t, err)
}
