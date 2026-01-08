package service

import (
	"context"
	"errors"

	"bizbundl/internal/constants"
	db "bizbundl/internal/db/sqlc/platform"
	"bizbundl/token"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// PlatformStore wrapper to match expected interface if needed, or just use Queries
// For MVP, directly use the generated Queries struct
type AuthService struct {
	store      *db.Queries
	tokenMaker token.Maker
}

func NewAuthService(store *db.Queries, tokenMaker token.Maker) *AuthService {
	return &AuthService{store: store, tokenMaker: tokenMaker}
}

// Helper for Module init
func NewPlatformQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
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
	hashed, err := hashPassword(password)
	if err != nil {
		return db.User{}, err
	}

	var firstName, lastName string
	firstName = fullName // Simplified split or passed from handler
	lastName = ""

	var user db.User
	user, err = s.store.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: hashed,
		FirstName:    firstName,
		LastName:     lastName,
		// No Phone in Platform User Schema
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
	token, _, err := s.tokenMaker.CreateToken(user.ID.String(), "owner", duration) // Hardcode role 'owner' for now
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
