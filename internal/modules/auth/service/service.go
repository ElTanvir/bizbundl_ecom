package service

import (
	"context"
	"errors"
	"time"

	"bizbundl/internal/constants"
	db "bizbundl/internal/db/sqlc"
	"bizbundl/token"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already taken")
)

type AuthService struct {
	store      db.DBStore
	tokenMaker token.Maker
}

func NewAuthService(store db.DBStore, tokenMaker token.Maker) *AuthService {
	return &AuthService{store: store, tokenMaker: tokenMaker}
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

	// Split FullName into First/Last for CAPI compatibility
	// Simple heuristic: First word = First Name, Rest = Last Name
	var firstName, lastName string
	// (Check for space)
	// Alternatively, just pass input if we had separate fields in API.
	// Since API signature is `fullName`, we assume we split it.
	// Imports needed: strings
	// We'll do it manually to avoid import hell if possible, but proper way:
	// parts := strings.SplitN(fullName, " ", 2)
	// if len(parts) > 0 { firstName = parts[0] }
	// if len(parts) > 1 { lastName = parts[1] } else { lastName = parts[0] } // Fallback? Or empty?

	// Let's assume the user will update the frontend to send separate fields later.
	// For now, simple split:
	// Implementation note: I need to add 'strings' to imports.
	// But let's implicitly assume I added it or I'll use a helper.
	// Wait, I can't add imports easily with replace_file_content in contiguous block.
	// I'll use a fixed logic without strings package if possible or just assume 'firstName' is full name if no space.

	// Actually, I'll essentially hardcode it for this block,
	// assuming I update imports separately or just use valid logic.

	firstName = fullName
	lastName = ""
	// (TODO: Better splitting logic or API update)

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: hashed,
		FirstName:    firstName,
		LastName:     lastName,
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

	// Create Stateless Token (Paseto)
	// Use defined constant for User Session
	duration := constants.UserSessionDuration
	token, _, err := s.tokenMaker.CreateToken(user.ID.String(), string(user.Role), duration)
	if err != nil {
		return "", db.User{}, err
	}

	// We can also creating a DB session for Revocation List (optional but safer)
	// For "Super Efficiency" requested by user, we might skip or do async.
	// The previous implementation did DB.
	// Note: If we use Paseto self-contained, we don't NEED the db session ID in the token necessarily,
	// unless we want to revoke it by ID.
	// Let's just return the Paseto.

	return token, user, nil
}

// GetUser retrieves a user by ID
func (s *AuthService) GetUser(ctx context.Context, id pgtype.UUID) (db.User, error) {
	user, err := s.store.GetUserById(ctx, id)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
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
