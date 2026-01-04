package auth

import (
	"context"
	"errors"
	"time"

	db "bizbundl/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already taken")
)

type AuthService struct {
	store db.DBStore
}

func NewAuthService(store db.DBStore) *AuthService {
	return &AuthService{store: store}
}

// hashPassword generates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// verifyPassword checks if the provided password matches the hash
func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Register creates a new user and returns the user object
func (s *AuthService) Register(ctx context.Context, email, password, fullName, phone string) (db.User, error) {
	// 1. Check if email exists (Optimistic handling via DB Unique constraint is better for race conditions)
	// But explicit check provides better error message immediately if desired.
	// We rely on DB constraint for atomic safety.

	hashed, err := hashPassword(password)
	if err != nil {
		return db.User{}, err
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: hashed,
		FullName:     fullName,
		Phone:        &phone,
		Role:         db.UserRoleCustomer,
	})
	if err != nil {
		// Check for specific PG Error (Unique Violation)
		// For MVP, just generic error or log it.
		return db.User{}, err
	}
	return user, nil
}

// Login verifies credentials and returns a Session Token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, db.User, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		// Avoid leaking if user exists or not, but for MVP standard error
		return "", db.User{}, ErrInvalidCredentials
	}

	if !verifyPassword(password, user.PasswordHash) {
		return "", db.User{}, ErrInvalidCredentials
	}

	// Create Session
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * 7 * time.Hour) // 7 Days

	_, err = s.store.CreateSession(ctx, db.CreateSessionParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return "", db.User{}, err
	}

	return token, user, nil
}

// VerifySession checks if value exists and is valid
func (s *AuthService) VerifySession(ctx context.Context, token string) (db.User, error) {
	session, err := s.store.GetSession(ctx, token)
	if err != nil {
		return db.User{}, errors.New("invalid session")
	}

	if session.ExpiresAt.Time.Before(time.Now()) {
		// Cleanup expired
		_ = s.store.DeleteSession(ctx, token)
		return db.User{}, errors.New("expired session")
	}

	// UserID is a pgtype.UUID, which might not be valid if it was a guest session
	// But Login always sets UserID.
	// If we support Guest Sessions later, we need to handle that.
	// The previous schema says UserID in Session is nullable (linked to user OR guest).
	// Wait, the schema says: `user_id UUID REFERENCES users(id)`. It is nullable.
	// Our Login function sets it.
	// If it's nil, we can't return a User.
	if !session.UserID.Valid {
		return db.User{}, errors.New("guest session")
	}

	user, err := s.store.GetUserById(ctx, session.UserID)
	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

// Logout deletes the session
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.store.DeleteSession(ctx, token)
}
